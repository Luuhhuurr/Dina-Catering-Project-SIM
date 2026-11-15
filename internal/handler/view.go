package handler

import (
	"sort"
	"time"

	"dinacatering/internal/store"
)

// HariProduksi adalah kumpulan tugas produksi untuk satu tanggal.
type HariProduksi struct {
	Tanggal time.Time
	Tugas   []store.ListTugasRentangTanggalRow
}

// kelompokkanPerHari mengelompokkan tugas produksi berdasarkan tanggal
// tenggatnya, lalu mengurutkannya dari yang paling dekat.
func kelompokkanPerHari(tugas []store.ListTugasRentangTanggalRow) []HariProduksi {
	group := map[string]*HariProduksi{}
	for _, t := range tugas {
		kunci := t.Deadline.Time.Format("2006-01-02")
		if _, ok := group[kunci]; !ok {
			group[kunci] = &HariProduksi{Tanggal: t.Deadline.Time}
		}
		group[kunci].Tugas = append(group[kunci].Tugas, t)
	}

	hasil := make([]HariProduksi, 0, len(group))
	for _, hari := range group {
		hasil = append(hasil, *hari)
	}
	sort.Slice(hasil, func(i, j int) bool {
		return hasil[i].Tanggal.Before(hasil[j].Tanggal)
	})
	return hasil
}
