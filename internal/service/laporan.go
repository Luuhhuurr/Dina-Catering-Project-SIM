package service

import (
	"context"
	"fmt"
	"time"

	"dinacatering/internal/store"
)

// LaporanLabaRugi merangkum hasil usaha sebuah periode.
type LaporanLabaRugi struct {
	Dari             time.Time
	Sampai           time.Time
	Pemasukan        int64
	PengeluaranBahan int64
	PengeluaranOps   int64
	Keuntungan       int64
}

// LaporanLabaRugiPeriode menghitung pemasukan, pengeluaran, dan keuntungan
// untuk rentang tanggal (aturan bisnis No.9).
func (s *Service) LaporanLabaRugiPeriode(ctx context.Context, dari, sampai time.Time) (LaporanLabaRugi, error) {
	dari = TanggalNormalize(dari)
	sampai = TanggalNormalize(sampai)

	pemasukan, err := s.Queries.TotalPenerimaanRentang(ctx, store.TotalPenerimaanRentangParams{
		Dari:   dari,
		Sampai: sampai,
	})
	if err != nil {
		return LaporanLabaRugi{}, fmt.Errorf("total pemasukan: %w", err)
	}
	bahan, err := s.Queries.TotalPembelianRentang(ctx, store.TotalPembelianRentangParams{
		Dari:   dari,
		Sampai: sampai,
	})
	if err != nil {
		return LaporanLabaRugi{}, fmt.Errorf("total pembelian bahan: %w", err)
	}
	ops, err := s.Queries.TotalBiayaRentang(ctx, store.TotalBiayaRentangParams{
		Dari:   dari,
		Sampai: sampai,
	})
	if err != nil {
		return LaporanLabaRugi{}, fmt.Errorf("total biaya operasional: %w", err)
	}

	return LaporanLabaRugi{
		Dari:             dari,
		Sampai:           sampai,
		Pemasukan:        pemasukan,
		PengeluaranBahan: bahan,
		PengeluaranOps:   ops,
		Keuntungan:       pemasukan - bahan - ops,
	}, nil
}

// MarginPesanan adalah indicator untung sebuah pesanan.
type MarginPesanan struct {
	PesananID     int32
	PelangganNama string
	PaketNama     string
	JumlahTamu    int32
	TotalHarga    int64
	HPPPerPorsi   int64
	TotalHPP      int64
	MarginPersen  float64
	Peringatan    bool
}

// MarginPerPesanan menghitung margin untung sebuah pesanan berdasarkan HPP
// paket pada tanggal acara.
func (s *Service) MarginPerPesanan(ctx context.Context, pesananID int32) (MarginPesanan, error) {
	pesanan, err := s.Queries.GetPesanan(ctx, pesananID)
	if err != nil {
		return MarginPesanan{}, fmt.Errorf("ambil pesanan: %w", err)
	}
	paket, err := s.Queries.GetPaket(ctx, pesanan.PaketMenuID)
	if err != nil {
		return MarginPesanan{}, fmt.Errorf("ambil paket: %w", err)
	}
	persen, err := s.PersenBiayaOperasional(ctx)
	if err != nil {
		return MarginPesanan{}, err
	}
	hpp, err := s.HitungHPPPaket(ctx, s.Queries, paket.ID, pesanan.TanggalAcara, persen)
	if err != nil {
		return MarginPesanan{}, err
	}

	hargaEfektif := HargaPerPorsiEfektif(paket.HargaPerPorsi, paket.HargaPerPorsiWeekend, pesanan.TanggalAcara)
	totalHPP := hpp.HPPPerPorsi * int64(pesanan.JumlahTamu)
	margin := MarginPaket(hargaEfektif, hpp.HPPPerPorsi)

	return MarginPesanan{
		PesananID:     pesanan.ID,
		PelangganNama: pesanan.PelangganNama,
		PaketNama:     pesanan.PaketNama,
		JumlahTamu:    pesanan.JumlahTamu,
		TotalHarga:    pesanan.TotalHarga,
		HPPPerPorsi:   hpp.HPPPerPorsi,
		TotalHPP:      totalHPP,
		MarginPersen:  margin,
		Peringatan:    MarginDiBawahMinimum(margin),
	}, nil
}

// MarginSeluruhPesanan menghitung margin untuk sekumpulan pesanan (dipakai
// pada laporan).
func (s *Service) MarginSeluruhPesanan(ctx context.Context, pesanan []store.ListPesananRow, t time.Time) ([]MarginPesanan, error) {
	persen, err := s.PersenBiayaOperasional(ctx)
	if err != nil {
		return nil, err
	}
	cache := map[int32]RingkasanHPPPaket{}

	hasil := make([]MarginPesanan, 0, len(pesanan))
	for _, p := range pesanan {
		hpp, ok := cache[p.PaketMenuID]
		if !ok {
			r, err := s.HitungHPPPaket(ctx, s.Queries, p.PaketMenuID, t, persen)
			if err != nil {
				return nil, err
			}
			hpp = r
			cache[p.PaketMenuID] = r
		}
		paket, err := s.Queries.GetPaket(ctx, p.PaketMenuID)
		if err != nil {
			return nil, fmt.Errorf("ambil paket: %w", err)
		}
		hargaEfektif := HargaPerPorsiEfektif(paket.HargaPerPorsi, paket.HargaPerPorsiWeekend, p.TanggalAcara)
		margin := MarginPaket(hargaEfektif, hpp.HPPPerPorsi)
		hasil = append(hasil, MarginPesanan{
			PesananID:     p.ID,
			PelangganNama: p.PelangganNama,
			PaketNama:     p.PaketNama,
			JumlahTamu:    p.JumlahTamu,
			TotalHarga:    p.TotalHarga,
			HPPPerPorsi:   hpp.HPPPerPorsi,
			TotalHPP:      hpp.HPPPerPorsi * int64(p.JumlahTamu),
			MarginPersen:  margin,
			Peringatan:    MarginDiBawahMinimum(margin),
		})
	}
	return hasil, nil
}
