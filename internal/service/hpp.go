package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"dinacatering/internal/model"
	"dinacatering/internal/store"
)

// KomponenHPP adalah satu baris resep beserta biaya per porsinya.
type KomponenHPP struct {
	BahanID        int32
	NamaBahan      string
	Satuan         string
	JumlahPerPorsi float64
	HargaSatuan    int64
	Biaya          int64
}

// RingkasanHPP adalah hasil perhitungan HPP untuk satu menu.
type RingkasanHPP struct {
	MenuID            int32
	MenuNama          string
	Komponen          []KomponenHPP
	HPPBahan          int64
	PersenOperasional float64
	HPPPerPorsi       int64
}

// RingkasanHPPPaket adalah hasil perhitungan HPP untuk satu paket menu.
type RingkasanHPPPaket struct {
	PaketID     int32
	PaketNama   string
	HPPPerPorsi int64
	Menu        []RingkasanHPP
}

// BiayaKomponen menghitung biaya satu baris resep untuk satu porsi.
func BiayaKomponen(jumlahPerPorsi float64, hargaSatuan int64) int64 {
	return int64(math.Round(jumlahPerPorsi * float64(hargaSatuan)))
}

// TotalHPPPerPorsi menerapkan proporsi biaya operasional pada HPP bahan.
//
// Asumsi bisnis: biaya operasional (default 10%) ditambahkan proporsional ke
// HPP bahan, lalu margin dihitung dari HPP total. Dengan cara ini rumus
// "harga jual minimum = HPP / 0.7" tetap konsisten.
func TotalHPPPerPorsi(hppBahan int64, persenOperasional float64) int64 {
	return int64(math.Round(float64(hppBahan) * (1 + persenOperasional/100)))
}

// MarginPaket menghitung margin untung sebuah paket dalam persen.
// Margin = (harga jual - HPP) / harga jual.
func MarginPaket(hargaJual, hpp int64) float64 {
	if hargaJual <= 0 {
		return 0
	}
	return float64(hargaJual-hpp) / float64(hargaJual) * 100
}

// HargaJualMinimum30Persen mengembalikan harga jual terendah agar margin
// minimal 30% (aturan bisnis No.8): HPP / 0.7.
//
// Toleransi kecil dipakai agar hasil bagi yang sebenarnya bulat (mis. HPP
// 15.400 -> 22.000) tidak ikut terbulatkan ke atas.
func HargaJualMinimum30Persen(hpp int64) int64 {
	if hpp <= 0 {
		return 0
	}
	return int64(math.Ceil(float64(hpp)/0.7 - 1e-9))
}

// MarginDiBawahMinimum menandai paket yang marginnya di bawah 30%.
func MarginDiBawahMinimum(margin float64) bool {
	return margin < model.MarginMinimumPersen
}

// hargaEfektifBahan mengembalikan harga satuan setiap bahan yang berlaku pada
// tanggal t. Bila bahan belum punya riwayat harga pada tanggal tersebut,
// dipakai harga terakhirnya.
func (s *Service) hargaEfektifBahan(ctx context.Context, q store.Querier, t time.Time) (map[int32]int64, error) {
	harga := map[int32]int64{}

	bahan, err := q.ListBahanBaku(ctx)
	if err != nil {
		return nil, fmt.Errorf("ambil bahan baku: %w", err)
	}
	for _, b := range bahan {
		harga[b.ID] = b.HargaTerakhir
	}

	riwayat, err := q.ListHargaEfektif(ctx, t)
	if err != nil {
		return nil, fmt.Errorf("ambil harga bahan efektif: %w", err)
	}
	for _, r := range riwayat {
		harga[r.BahanBakuID] = r.HargaSatuan
	}
	return harga, nil
}

// HitungHPPMenu menghitung HPP per porsi untuk satu menu pada tanggal t.
func (s *Service) HitungHPPMenu(ctx context.Context, q store.Querier, menuID int32, t time.Time, persenOperasional float64) (RingkasanHPP, error) {
	menu, err := q.GetMenu(ctx, menuID)
	if err != nil {
		return RingkasanHPP{}, fmt.Errorf("ambil menu: %w", err)
	}
	return s.hitungHPPMenuResep(ctx, q, menu, t, persenOperasional)
}

func (s *Service) hitungHPPMenuResep(ctx context.Context, q store.Querier, menu store.Menu, t time.Time, persenOperasional float64) (RingkasanHPP, error) {
	resep, err := q.ListResepByMenu(ctx, menu.ID)
	if err != nil {
		return RingkasanHPP{}, fmt.Errorf("ambil resep menu %d: %w", menu.ID, err)
	}
	harga, err := s.hargaEfektifBahan(ctx, q, t)
	if err != nil {
		return RingkasanHPP{}, err
	}

	ringkasan := RingkasanHPP{
		MenuID:            menu.ID,
		MenuNama:          menu.Nama,
		PersenOperasional: persenOperasional,
		Komponen:          make([]KomponenHPP, 0, len(resep)),
	}
	var hppBahan int64
	for _, r := range resep {
		hargaSatuan := harga[r.BahanBakuID]
		biaya := BiayaKomponen(r.JumlahPerPorsi, hargaSatuan)
		hppBahan += biaya
		ringkasan.Komponen = append(ringkasan.Komponen, KomponenHPP{
			BahanID:        r.BahanBakuID,
			NamaBahan:      r.BahanNama,
			Satuan:         r.BahanSatuan,
			JumlahPerPorsi: r.JumlahPerPorsi,
			HargaSatuan:    hargaSatuan,
			Biaya:          biaya,
		})
	}
	ringkasan.HPPBahan = hppBahan
	ringkasan.HPPPerPorsi = TotalHPPPerPorsi(hppBahan, persenOperasional)
	return ringkasan, nil
}

// HitungHPPPaket menghitung HPP per porsi untuk satu paket pada tanggal t.
func (s *Service) HitungHPPPaket(ctx context.Context, q store.Querier, paketID int32, t time.Time, persenOperasional float64) (RingkasanHPPPaket, error) {
	paket, err := q.GetPaket(ctx, paketID)
	if err != nil {
		return RingkasanHPPPaket{}, fmt.Errorf("ambil paket: %w", err)
	}
	return s.hitungHPPPaketDetail(ctx, q, paket.ID, paket.Nama, t, persenOperasional)
}

func (s *Service) hitungHPPPaketDetail(ctx context.Context, q store.Querier, paketID int32, paketNama string, t time.Time, persenOperasional float64) (RingkasanHPPPaket, error) {
	item, err := q.ListPaketItem(ctx, paketID)
	if err != nil {
		return RingkasanHPPPaket{}, fmt.Errorf("ambil isi paket %d: %w", paketID, err)
	}

	ringkasan := RingkasanHPPPaket{
		PaketID:   paketID,
		PaketNama: paketNama,
		Menu:      make([]RingkasanHPP, 0, len(item)),
	}
	for _, it := range item {
		menu, err := q.GetMenu(ctx, it.MenuID)
		if err != nil {
			return RingkasanHPPPaket{}, fmt.Errorf("ambil menu %d: %w", it.MenuID, err)
		}
		hpp, err := s.hitungHPPMenuResep(ctx, q, menu, t, persenOperasional)
		if err != nil {
			return RingkasanHPPPaket{}, err
		}
		ringkasan.Menu = append(ringkasan.Menu, hpp)
		ringkasan.HPPPerPorsi += hpp.HPPPerPorsi
	}
	return ringkasan, nil
}

// RingkasanPaket menggabungkan data paket, HPP, margin, dan peringatan
// margin di bawah 30% untuk ditampilkan pada daftar paket.
type RingkasanPaket struct {
	PaketID              int32
	Nama                 string
	HargaPerPorsi        int64
	HargaPerPorsiWeekend *int64
	TahunBerlaku         int32
	Deskripsi            string
	JumlahMenu           int64
	HPPPerPorsi          int64
	MarginPersen         float64
	HargaJualMinimum     int64
	HargaJualEfektif     int64
	Menu                 []RingkasanHPP
	PeringatanMargin     bool
}

// DaftarPaket menghitung HPP dan margin untuk seluruh paket pada tanggal t.
func (s *Service) DaftarPaket(ctx context.Context, t time.Time) ([]RingkasanPaket, error) {
	persen, err := s.PersenBiayaOperasional(ctx)
	if err != nil {
		return nil, err
	}
	paket, err := s.Queries.ListPaket(ctx)
	if err != nil {
		return nil, fmt.Errorf("ambil daftar paket: %w", err)
	}

	hasil := make([]RingkasanPaket, 0, len(paket))
	for _, p := range paket {
		r, err := s.hitungHPPPaketDetail(ctx, s.Queries, p.ID, p.Nama, t, persen)
		if err != nil {
			return nil, err
		}
		hargaEfektif := HargaPerPorsiEfektif(p.HargaPerPorsi, p.HargaPerPorsiWeekend, t)
		margin := MarginPaket(hargaEfektif, r.HPPPerPorsi)
		hasil = append(hasil, RingkasanPaket{
			PaketID:              p.ID,
			Nama:                 p.Nama,
			HargaPerPorsi:        p.HargaPerPorsi,
			HargaPerPorsiWeekend: p.HargaPerPorsiWeekend,
			TahunBerlaku:         p.TahunBerlaku,
			Deskripsi:            p.Deskripsi,
			JumlahMenu:           p.JumlahMenu,
			HPPPerPorsi:          r.HPPPerPorsi,
			MarginPersen:         margin,
			HargaJualMinimum:     HargaJualMinimum30Persen(r.HPPPerPorsi),
			HargaJualEfektif:     hargaEfektif,
			Menu:                 r.Menu,
			PeringatanMargin:     MarginDiBawahMinimum(margin),
		})
	}
	return hasil, nil
}

// PaketDenganMarginDiBawahMinimum mengembalikan paket yang marginnya di bawah
// 30% — dipakai untuk peringatan di dashboard dan laporan.
func (s *Service) PaketDenganMarginDiBawahMinimum(ctx context.Context, t time.Time) ([]RingkasanPaket, error) {
	semua, err := s.DaftarPaket(ctx, t)
	if err != nil {
		return nil, err
	}
	var bermasalah []RingkasanPaket
	for _, p := range semua {
		if p.PeringatanMargin {
			bermasalah = append(bermasalah, p)
		}
	}
	return bermasalah, nil
}
