package service

import (
	"testing"
	"time"

	"dinacatering/internal/model"
)

func tanggal(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func TestJadwalPembayaran(t *testing.T) {
	// acara 15 Desember 2026 -> DP H-30 = 15 November 2026, pelunasan H-7 = 8 Desember 2026
	dp, lunas := JadwalPembayaran(tanggal(2026, time.December, 15))
	if !dp.Equal(tanggal(2026, time.November, 15)) {
		t.Fatalf("tenggat DP = %v, ingin 2026-11-15", dp)
	}
	if !lunas.Equal(tanggal(2026, time.December, 8)) {
		t.Fatalf("tenggat pelunasan = %v, ingin 2026-12-08", lunas)
	}
}

func TestDPAmount(t *testing.T) {
	// DP 50% dari total pesanan
	if got := DPAmount(3500000); got != 1750000 {
		t.Fatalf("DPAmount(3500000) = %d, ingin 1750000", got)
	}
}

func TestStatusPembayaran(t *testing.T) {
	if got := StatusPembayaran(100000, 0); got != model.StatusBelumDP {
		t.Fatalf("status tanpa pembayaran = %q, ingin %q", got, model.StatusBelumDP)
	}
	if got := StatusPembayaran(100000, 50000); got != model.StatusDPTerbayar {
		t.Fatalf("status setelah DP = %q, ingin %q", got, model.StatusDPTerbayar)
	}
	if got := StatusPembayaran(100000, 100000); got != model.StatusLunas {
		t.Fatalf("status setelah pelunasan = %q, ingin %q", got, model.StatusLunas)
	}
}

func TestStatusJatuhTempo(t *testing.T) {
	acara := tanggal(2026, time.December, 15)
	hariIni := tanggal(2026, time.November, 10)

	// belum ada pembayaran, hari ini masih sebelum H-30 (15 November)
	if got := StatusJatuhTempo(acara, 100000, 0, hariIni); got != "dp_jatuh_tempo" {
		t.Fatalf("status = %q, ingin dp_jatuh_tempo", got)
	}
	// sudah lewat H-30 tanpa DP
	if got := StatusJatuhTempo(acara, 100000, 0, tanggal(2026, time.December, 1)); got != "dp_terlambat" {
		t.Fatalf("status = %q, ingin dp_terlambat", got)
	}
	// DP sudah dibayar, tinggal pelunasan
	if got := StatusJatuhTempo(acara, 100000, 50000, hariIni); got != "pelunasan_jatuh_tempo" {
		t.Fatalf("status = %q, ingin pelunasan_jatuh_tempo", got)
	}
	// sudah lunas
	if got := StatusJatuhTempo(acara, 100000, 100000, hariIni); got != "lunas" {
		t.Fatalf("status = %q, ingin lunas", got)
	}
}

func TestIsWeekend(t *testing.T) {
	// 14 November 2026 adalah Sabtu, 15 November 2026 Minggu, 16 November 2026 Senin
	if !IsWeekend(tanggal(2026, time.November, 14)) {
		t.Fatal("Sabtu seharusnya weekend")
	}
	if !IsWeekend(tanggal(2026, time.November, 15)) {
		t.Fatal("Minggu seharusnya weekend")
	}
	if IsWeekend(tanggal(2026, time.November, 16)) {
		t.Fatal("Senin bukan weekend")
	}
}

func TestHargaPerPorsiEfektif(t *testing.T) {
	weekend := int64(38000)
	biasa := int64(35000)

	// hari biasa pakai harga biasa
	if got := HargaPerPorsiEfektif(biasa, &weekend, tanggal(2026, time.November, 16)); got != 35000 {
		t.Fatalf("harga hari biasa = %d, ingin 35000", got)
	}
	// hari Sabtu pakai harga weekend
	if got := HargaPerPorsiEfektif(biasa, &weekend, tanggal(2026, time.November, 14)); got != 38000 {
		t.Fatalf("harga weekend = %d, ingin 38000", got)
	}
	// tanpa harga weekend memakai harga biasa
	if got := HargaPerPorsiEfektif(biasa, nil, tanggal(2026, time.November, 14)); got != 35000 {
		t.Fatalf("harga weekend tanpa kolom = %d, ingin 35000", got)
	}
}

func TestHitungTotalPesanan(t *testing.T) {
	// 60 tamu x 38.000 (harga weekend paket Standar) = 2.280.000
	if got := HitungTotalPesanan(60, 38000); got != 2280000 {
		t.Fatalf("total pesanan = %d, ingin 2280000", got)
	}
	if got := HitungTotalPesanan(0, 38000); got != 0 {
		t.Fatalf("total pesanan tanpa tamu = %d, ingin 0", got)
	}
}
