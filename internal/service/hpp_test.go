package service

import (
	"math"
	"testing"
)

func TestBiayaKomponen(t *testing.T) {
	// 0,15 kg sayur dengan harga 8.000/kg = 1.200 per porsi
	if got := BiayaKomponen(0.15, 8000); got != 1200 {
		t.Fatalf("BiayaKomponen(0.15, 8000) = %d, ingin 1200", got)
	}
	// 0,125 kg beras dengan harga 16.000/kg = 2.000 per porsi
	if got := BiayaKomponen(0.125, 16000); got != 2000 {
		t.Fatalf("BiayaKomponen(0.125, 16000) = %d, ingin 2000", got)
	}
	// pembulatan ke rupiah terdekat
	if got := BiayaKomponen(0.033, 12000); got != 396 {
		t.Fatalf("BiayaKomponen(0.033, 12000) = %d, ingin 396", got)
	}
}

func TestTotalHPPPerPorsi(t *testing.T) {
	// HPP bahan 21.000 + 10% biaya operasional = 23.100
	if got := TotalHPPPerPorsi(21000, 10); got != 23100 {
		t.Fatalf("TotalHPPPerPorsi(21000, 10) = %d, ingin 23100", got)
	}
	// tanpa biaya operasional HPP tetap sama
	if got := TotalHPPPerPorsi(21000, 0); got != 21000 {
		t.Fatalf("TotalHPPPerPorsi(21000, 0) = %d, ingin 21000", got)
	}
}

func TestMarginPaket(t *testing.T) {
	// harga jual 35.000, HPP 15.400 -> margin 56%
	margin := MarginPaket(35000, 15400)
	if math.Abs(margin-56) > 0.001 {
		t.Fatalf("MarginPaket(35000, 15400) = %v, ingin 56", margin)
	}
	// HPP lebih besar dari harga jual -> margin negatif
	if got := MarginPaket(30000, 33000); got >= 0 {
		t.Fatalf("MarginPaket(30000, 33000) = %v, ingin negatif", got)
	}
	// harga jual tidak valid -> 0
	if got := MarginPaket(0, 1000); got != 0 {
		t.Fatalf("MarginPaket(0, 1000) = %v, ingin 0", got)
	}
}

func TestHargaJualMinimum30Persen(t *testing.T) {
	// HPP 15.400 -> 15.400 / 0.7 = 22.000
	if got := HargaJualMinimum30Persen(15400); got != 22000 {
		t.Fatalf("HargaJualMinimum30Persen(15400) = %d, ingin 22000", got)
	}
	// hasil harus dibulatkan ke atas agar margin tidak pernah di bawah 30%
	hpp := int64(23067)
	minimum := HargaJualMinimum30Persen(hpp)
	if MarginPaket(minimum, hpp) < 29.99 {
		t.Fatalf("harga minimum %d memberi margin %.4f%%, ingin >= 30%%", minimum, MarginPaket(minimum, hpp))
	}
}

func TestMarginDiBawahMinimum(t *testing.T) {
	if !MarginDiBawahMinimum(29.9) {
		t.Fatal("margin 29.9% seharusnya berstatus di bawah minimum")
	}
	if MarginDiBawahMinimum(30) {
		t.Fatal("margin 30% tidak seharusnya berstatus di bawah minimum")
	}
	if MarginDiBawahMinimum(56) {
		t.Fatal("margin 56% tidak seharusnya berstatus di bawah minimum")
	}
}

func TestHitungHPPContohPaketPremium(t *testing.T) {
	// Angka berikut meniru data awal paket Premium (Nasi + Rendang + Sayur + Buah):
	//   Nasi Putih     2.670
	//   Rendang Daging 20.970
	//   Sayur Sop      2.110
	//   Buah Potong    2.250
	hppBahan := int64(2670 + 20970 + 2110 + 2250)
	if hppBahan != 28000 {
		t.Fatalf("HPP bahan = %d, ingin 28000", hppBahan)
	}
	hpp := TotalHPPPerPorsi(hppBahan, 10)
	if hpp != 30800 {
		t.Fatalf("HPP per porsi = %d, ingin 30800", hpp)
	}

	margin := MarginPaket(40000, hpp)
	if margin >= 30 {
		t.Fatalf("margin paket Premium = %.2f%%, ingin di bawah 30%%", margin)
	}
	if !MarginDiBawahMinimum(margin) {
		t.Fatal("paket Premium seharusnya memunculkan peringatan margin")
	}
	if min := HargaJualMinimum30Persen(hpp); min != 44000 {
		t.Fatalf("harga jual minimum = %d, ingin 44000", min)
	}

	// Harga akhir pekan 43.000 pun masih di bawah margin minimum, sehingga
	// peringatan tetap muncul baik hari kerja maupun akhir pekan.
	marginAkhirPekan := MarginPaket(43000, hpp)
	if !MarginDiBawahMinimum(marginAkhirPekan) {
		t.Fatalf("margin akhir pekan = %.2f%%, ingin di bawah 30%%", marginAkhirPekan)
	}
}
