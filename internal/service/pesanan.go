package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"dinacatering/internal/model"
	"dinacatering/internal/store"
)

// PesananInput adalah data formulir pesanan yang diterima dari handler.
type PesananInput struct {
	PelangganID   int32
	NamaPaket     string
	JumlahTamu    int
	TanggalAcara  time.Time
	Lokasi        string
	SumberPesanan string
	Catatan       string
	MasakDiLokasi bool
}

// BuatPesanan membuat pesanan baru. Total harga dan tenggat DP/pelunasan
// dihitung otomatis dari tanggal acara dan harga paket yang berlaku.
func (s *Service) BuatPesanan(ctx context.Context, in PesananInput) (int32, error) {
	if in.PelangganID <= 0 {
		return 0, errors.New("pelanggan wajib dipilih")
	}
	if in.NamaPaket == "" {
		return 0, errors.New("paket menu wajib dipilih")
	}
	if in.JumlahTamu <= 0 {
		return 0, errors.New("jumlah tamu harus lebih dari nol")
	}
	if in.SumberPesanan == "" {
		in.SumberPesanan = model.SumberLainnya
	}

	tanggalAcara := TanggalNormalize(in.TanggalAcara)
	paket, err := s.ResolvePaketUntukTahun(ctx, in.NamaPaket, tanggalAcara)
	if err != nil {
		return 0, err
	}

	hargaEfektif := HargaPerPorsiEfektif(paket.HargaPerPorsi, paket.HargaPerPorsiWeekend, tanggalAcara)
	total := HitungTotalPesanan(in.JumlahTamu, hargaEfektif)
	dp, pelunasan := JadwalPembayaran(tanggalAcara)

	var id int32
	err = s.withTx(ctx, func(q store.Querier) error {
		pesananID, err := q.CreatePesanan(ctx, store.CreatePesananParams{
			PelangganID:        in.PelangganID,
			PaketMenuID:        paket.ID,
			JumlahTamu:         int32(in.JumlahTamu),
			TanggalAcara:       tanggalAcara,
			Lokasi:             in.Lokasi,
			SumberPesanan:      in.SumberPesanan,
			TotalHarga:         total,
			Status:             model.StatusBaru,
			Catatan:            in.Catatan,
			MasakDiLokasi:      in.MasakDiLokasi,
			TglJatuhTempoDp:    dp,
			TglJatuhTempoLunas: pelunasan,
		})
		if err != nil {
			return fmt.Errorf("simpan pesanan: %w", err)
		}
		id = pesananID
		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateStatusPesanan mengubah status alur pesanan. Saat pesanan masuk status
// produksi, stok bahan dikurangi sesuai resep, dan tugas produksi default
// dibuat bila belum ada.
func (s *Service) UpdateStatusPesanan(ctx context.Context, pesananID int32, statusBaru string) error {
	pesanan, err := s.Queries.GetPesanan(ctx, pesananID)
	if err != nil {
		return fmt.Errorf("ambil pesanan: %w", err)
	}
	if pesanan.Status == statusBaru {
		return nil
	}

	return s.withTx(ctx, func(q store.Querier) error {
		if err := q.UpdateStatusPesanan(ctx, store.UpdateStatusPesananParams{
			ID:     pesananID,
			Status: statusBaru,
		}); err != nil {
			return fmt.Errorf("update status pesanan: %w", err)
		}

		if statusBaru == model.StatusProduksi && pesanan.Status != model.StatusProduksi {
			if err := s.kurangiStokUntukProduksi(ctx, q, pesanan); err != nil {
				return err
			}
			if err := s.buatTugasDefault(ctx, q, pesanan); err != nil {
				return err
			}
		}
		return nil
	})
}

// buatTugasDefault membuat daftar tugas produksi standar untuk hari acara:
// persiapan H-1, masak H-H, dan pengiriman (aturan bisnis No.7).
func (s *Service) buatTugasDefault(ctx context.Context, q store.Querier, pesanan store.GetPesananRow) error {
	existing, err := q.ListTugasByPesanan(ctx, pesanan.ID)
	if err != nil {
		return fmt.Errorf("ambil tugas produksi: %w", err)
	}
	if len(existing) > 0 {
		return nil
	}

	acara := TanggalNormalize(pesanan.TanggalAcara)
	persiapan := acara.AddDate(0, 0, -1)

	tugas := []struct {
		nama     string
		deadline time.Time
	}{
		{"Persiapan bahan & alat", persiapan},
		{"Masak seluruh menu", acara},
		{"Plating & cek porsi", acara},
		{"Pengiriman ke lokasi", acara},
	}
	for _, t := range tugas {
		_, err := q.InsertTugas(ctx, store.InsertTugasParams{
			PesananID: pesanan.ID,
			NamaTugas: t.nama,
			Pj:        "",
			Deadline:  pgDate(t.deadline),
			Status:    "belum",
		})
		if err != nil {
			return fmt.Errorf("buat tugas produksi: %w", err)
		}
	}
	return nil
}
