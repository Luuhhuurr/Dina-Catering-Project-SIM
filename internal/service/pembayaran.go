package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"dinacatering/internal/model"
	"dinacatering/internal/store"
)

// JadwalPembayaran menghitung tenggat DP (H-30) dan pelunasan (H-7) dari
// tanggal acara (aturan bisnis No.3).
func JadwalPembayaran(tanggalAcara time.Time) (dp time.Time, pelunasan time.Time) {
	acara := TanggalNormalize(tanggalAcara)
	return acara.AddDate(0, 0, -30), acara.AddDate(0, 0, -7)
}

// DPAmount adalah nominal DP 50% dari total pesanan.
func DPAmount(totalHarga int64) int64 {
	return totalHarga / 2
}

// StatusPembayaran menurunkan status pembayaran dari total yang sudah dibayar.
func StatusPembayaran(totalHarga, totalDibayar int64) string {
	switch {
	case totalHarga > 0 && totalDibayar >= totalHarga:
		return model.StatusLunas
	case totalDibayar > 0:
		return model.StatusDPTerbayar
	default:
		return model.StatusBelumDP
	}
}

// StatusJatuhTempo menandai posisi pembayaran terhadap tenggatnya sehingga
// owner bisa melihat pesanan yang mendesak (H-30 dan H-7).
func StatusJatuhTempo(tanggalAcara time.Time, totalHarga, totalDibayar int64, hariIni time.Time) string {
	if StatusPembayaran(totalHarga, totalDibayar) == model.StatusLunas {
		return "lunas"
	}
	dp, pelunasan := JadwalPembayaran(tanggalAcara)
	hariIni = TanggalNormalize(hariIni)

	switch StatusPembayaran(totalHarga, totalDibayar) {
	case model.StatusBelumDP:
		if !hariIni.After(dp) {
			return "dp_jatuh_tempo"
		}
		return "dp_terlambat"
	default:
		if !hariIni.After(pelunasan) {
			return "pelunasan_jatuh_tempo"
		}
		return "pelunasan_terlambat"
	}
}

// LabelStatusJatuhTempo memberi keterangan berbahasa manusia untuk daftar
// "menunggu pembayaran".
func LabelStatusJatuhTempo(status string) string {
	switch status {
	case "dp_jatuh_tempo":
		return "DP jatuh tempo"
	case "dp_terlambat":
		return "DP terlambat"
	case "pelunasan_jatuh_tempo":
		return "Pelunasan jatuh tempo"
	case "pelunasan_terlambat":
		return "Pelunasan terlambat"
	default:
		return "Lunas"
	}
}

// Pembayaran meringkas kondisi pembayaran sebuah pesanan.
type Pembayaran struct {
	Total   int64
	Status  string
	BelumDP bool
	SudahDP bool
	Lunas   bool
}

// RingkasanPembayaran menghitung status pembayaran sebuah pesanan.
func (s *Service) RingkasanPembayaran(ctx context.Context, pesananID int32, totalHarga int64) (Pembayaran, error) {
	total, err := s.Queries.TotalPembayaranPesanan(ctx, pesananID)
	if err != nil {
		return Pembayaran{}, fmt.Errorf("total pembayaran: %w", err)
	}
	status := StatusPembayaran(totalHarga, total)
	return Pembayaran{
		Total:   total,
		Status:  status,
		BelumDP: status == model.StatusBelumDP,
		SudahDP: status == model.StatusDPTerbayar,
		Lunas:   status == model.StatusLunas,
	}, nil
}

// CatatPembayaran menyimpan pembayaran pesanan.
func (s *Service) CatatPembayaran(ctx context.Context, pesananID int32, jumlah int64, tipe, metode string, tanggal time.Time) error {
	if jumlah <= 0 {
		return errors.New("nominal pembayaran harus lebih dari nol")
	}
	pesanan, err := s.Queries.GetPesanan(ctx, pesananID)
	if err != nil {
		return fmt.Errorf("ambil pesanan: %w", err)
	}
	if pesanan.Status == model.StatusBatal {
		return errors.New("pesanan sudah dibatalkan, pembayaran tidak bisa dicatat")
	}

	return s.withTx(ctx, func(q store.Querier) error {
		_, err := q.CreatePembayaran(ctx, store.CreatePembayaranParams{
			PesananID: pesananID,
			Jumlah:    jumlah,
			Tipe:      tipe,
			Tanggal:   TanggalNormalize(tanggal),
			Metode:    metode,
		})
		if err != nil {
			return fmt.Errorf("simpan pembayaran: %w", err)
		}

		// DP yang sudah masuk berarti pesanan pasti jalan: status naik dari
		// "baru" ke "final_meeting" dan jadwal produksi untuk hari acara
		// langsung dibuat. Pembuatan tugas bersifat idempoten, jadi tidak
		// menggandakan tugas saat status pesanan nanti berubah menjadi "produksi".
		if tipe == model.TipeDP && pesanan.Status == model.StatusBaru {
			if err := q.UpdateStatusPesanan(ctx, store.UpdateStatusPesananParams{
				ID:     pesananID,
				Status: model.StatusFinalMeeting,
			}); err != nil {
				return fmt.Errorf("update status pesanan: %w", err)
			}
			if err := s.buatTugasDefault(ctx, q, pesanan); err != nil {
				return err
			}
		}
		return nil
	})
}
