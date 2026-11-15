package service

import (
	"context"
	"fmt"
	"time"

	"dinacatering/internal/store"
)

// TanggalNormalize membuang bagian jam/menit/detik dari sebuah waktu tanggal.
func TanggalNormalize(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// hariIni mengembalikan tanggal hari ini tanpa jam, dipakai untuk membandingkan
// tenggat pembayaran.
func hariIni() time.Time {
	return TanggalNormalize(time.Now())
}

// IsWeekend menandai Sabtu dan Minggu — hari acara yang tarifnya bisa berbeda
// (aturan bisnis No.2).
func IsWeekend(t time.Time) bool {
	wd := TanggalNormalize(t).Weekday()
	return wd == time.Saturday || wd == time.Sunday
}

// HargaPerPorsiEfektif mengembalikan harga per porsi yang berlaku untuk paket
// pada tanggal acara: pakai harga akhir pekan bila ada, selain itu harga biasa.
func HargaPerPorsiEfektif(hargaPerPorsi int64, hargaPerPorsiWeekend *int64, tanggalAcara time.Time) int64 {
	if IsWeekend(tanggalAcara) && hargaPerPorsiWeekend != nil {
		return *hargaPerPorsiWeekend
	}
	return hargaPerPorsi
}

// ResolvePaketUntukTahun mencari baris harga paket yang berlaku untuk tahun
// acara: baris dengan tahun_berlaku terdekat yang tidak lebih besar dari tahun
// acara. Bila belum ada, dikembalikan ErrPaketTidakDitemukan.
func (s *Service) ResolvePaketUntukTahun(ctx context.Context, namaPaket string, tanggalAcara time.Time) (store.PaketMenu, error) {
	paket, err := s.Queries.GetPaketEfektif(ctx, store.GetPaketEfektifParams{
		Nama:  namaPaket,
		Tahun: int32(tanggalAcara.Year()),
	})
	if err != nil {
		return store.PaketMenu{}, fmt.Errorf("%w: %s tahun %d", ErrPaketTidakDitemukan, namaPaket, tanggalAcara.Year())
	}
	return paket, nil
}

// HitungTotalPesanan menghitung total harga pesanan = jumlah tamu x harga per
// porsi yang berlaku pada tanggal acara.
func HitungTotalPesanan(jumlahTamu int, hargaPerPorsi int64) int64 {
	if jumlahTamu <= 0 {
		return 0
	}
	return int64(jumlahTamu) * hargaPerPorsi
}

// DaftarPaketPerTahun mengembalikan paket yang bisa dipilih untuk tanggal acara
// tertentu (hanya paket yang sudah punya harga di tahun yang relevan).
func (s *Service) DaftarPaketPerTahun(ctx context.Context, tanggalAcara time.Time) ([]store.PaketMenu, error) {
	paket, err := s.Queries.ListPaketPerTahun(ctx, int32(tanggalAcara.Year()))
	if err != nil {
		return nil, fmt.Errorf("ambil paket per tahun: %w", err)
	}
	return paket, nil
}
