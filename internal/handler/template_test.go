package handler

import (
	"bytes"
	"testing"
	"time"

	"dinacatering/internal/model"
	"dinacatering/internal/service"
	"dinacatering/internal/store"
)

// TestSemuaTemplateDapatDiParse memastikan setiap halaman dan partial
// merupakan template HTML yang valid.
func TestSemuaTemplateDapatDiParse(t *testing.T) {
	templates, err := buildTemplates()
	if err != nil {
		t.Fatalf("gagal mem-parse template: %v", err)
	}

	halaman := []string{
		"login", "dashboard", "pelanggan", "bahan", "pembelian",
		"menu", "paket", "hpp", "pesanan", "produksi", "laporan",
		"biaya", "pengaturan",
	}
	for _, nama := range halaman {
		if templates[nama] == nil {
			t.Fatalf("template %s tidak ditemukan", nama)
		}
	}
}

// TestRenderHalamanKosong memanggil setiap halaman dengan data contoh sehingga
// salah satu nama field yang salah ketik langsung ketahuan.
func TestRenderHalamanKosong(t *testing.T) {
	templates, err := buildTemplates()
	if err != nil {
		t.Fatalf("gagal mem-parse template: %v", err)
	}

	hari := time.Date(2026, time.November, 15, 0, 0, 0, 0, time.Local)
	konten := map[string]any{
		"HariIni":           hari,
		"Tanggal":           "2026-11-15",
		"Dari":              hari,
		"Sampai":            hari,
		"Bulan":             "2026-11",
		"Mulai":             hari,
		"Selesai":           hari,
		"Persen":            10.0,
		"Cari":              "",
		"Pesan":             "",
		"Status":            "",
		"Sumber":            "",
		"JumlahPelanggan":   5,
		"StatusCounts":      []store.CountPesananByStatusRow{},
		"SumberCounts":      []store.CountPesananBySumberRow{},
		"StatusPesanan":     model.StatusPesananList,
		"SumberPesanan":     model.SumberPesananList,
		"Metode":            model.MetodePembayaranList,
		"Daftar":            []store.ListPelangganRow{},
		"Bahan":             []store.BahanBaku{},
		"Menu":              []store.ListMenuRow{},
		"Paket":             []store.ListPaketRow{},
		"Pembelian":         []store.ListPembelianRow{},
		"Biaya":             []store.BiayaOperasional{},
		"Pengaturan":        []store.Pengaturan{},
		"MenungguBayar":     []service.BarisPesanan{},
		"ProduksiMendatang": []store.ListPesananRentangTanggalRow{},
		"PeringatanMargin":  []service.RingkasanPaket{},
		"Margin":            []service.MarginPesanan{},
		"Papan":             []HariProduksi{},
		"LabaRugi":          service.LaporanLabaRugi{Dari: hari, Sampai: hari},
		"Total":             int64(0),
		"NilaiDP":           int64(500000),
		"NilaiPelunasan":    int64(500000),
		"Detail":            false,
		"Kosong":            true,
	}

	for nama := range templates {
		var buf bytes.Buffer
		data := map[string]any{
			"Judul":    "Uji",
			"Menu":     nama,
			"Username": "admin",
			"Content":  konten,
		}
		if err := templates[nama].ExecuteTemplate(&buf, "base", data); err != nil {
			t.Errorf("halaman %s gagal dirender: %v", nama, err)
		}
	}
}
