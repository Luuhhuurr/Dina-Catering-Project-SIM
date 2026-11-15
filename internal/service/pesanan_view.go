package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"dinacatering/internal/model"
	"dinacatering/internal/store"
)

// BarisPesanan adalah tampilan pesanan yang sudah dilengkapi status
// pembayaran, sisa tagihan, dan keterangan tenggat.
type BarisPesanan struct {
	ID                 int32
	PelangganNama      string
	PaketNama          string
	JumlahTamu         int32
	TanggalAcara       time.Time
	Lokasi             string
	SumberPesanan      string
	TotalHarga         int64
	Status             string
	Catatan            string
	MasakDiLokasi      bool
	TglJatuhTempoDP    time.Time
	TglJatuhTempoLunas time.Time
	TotalDibayar       int64
	StatusPembayaran   string
	LabelJatuhTempo    string
	SisaTagihan        int64
}

// filterKosong mengubah string kosong menjadi nil, yang berarti "semua" pada
// query berparameter nullable.
func filterKosong(nilai string) *string {
	if nilai == "" {
		return nil
	}
	return &nilai
}

// DaftarPesanan mengembalikan pesanan yang sudah dilengkapi status pembayaran.
// Filter status dan sumber yang kosong berarti "semua".
func (s *Service) DaftarPesanan(ctx context.Context, status, sumber string, dari, sampai time.Time) ([]BarisPesanan, error) {
	rows, err := s.Queries.ListPesanan(ctx, store.ListPesananParams{
		Status: filterKosong(status),
		Sumber: filterKosong(sumber),
		Dari:   dari,
		Sampai: sampai,
	})
	if err != nil {
		return nil, fmt.Errorf("ambil daftar pesanan: %w", err)
	}

	hasil := make([]BarisPesanan, 0, len(rows))
	for _, row := range rows {
		hasil = append(hasil, s.lengkapiPesanan(row))
	}
	return hasil, nil
}

// lengkapiPesanan menghitung field turunan untuk satu pesanan.
func (s *Service) lengkapiPesanan(row store.ListPesananRow) BarisPesanan {
	sisa := row.TotalHarga - row.TotalBayar
	if sisa < 0 {
		sisa = 0
	}
	return BarisPesanan{
		ID:                 row.ID,
		PelangganNama:      row.PelangganNama,
		PaketNama:          row.PaketNama,
		JumlahTamu:         row.JumlahTamu,
		TanggalAcara:       row.TanggalAcara,
		Lokasi:             row.Lokasi,
		SumberPesanan:      row.SumberPesanan,
		TotalHarga:         row.TotalHarga,
		Status:             row.Status,
		Catatan:            row.Catatan,
		MasakDiLokasi:      row.MasakDiLokasi,
		TglJatuhTempoDP:    row.TglJatuhTempoDp,
		TglJatuhTempoLunas: row.TglJatuhTempoLunas,
		TotalDibayar:       row.TotalBayar,
		StatusPembayaran:   StatusPembayaran(row.TotalHarga, row.TotalBayar),
		LabelJatuhTempo:    LabelStatusJatuhTempo(StatusJatuhTempo(row.TanggalAcara, row.TotalHarga, row.TotalBayar, hariIni())),
		SisaTagihan:        sisa,
	}
}

// PesananMenungguPembayaran mengembalikan pesanan yang belum lunas dan
// tenggatnya sudah dekat, diurutkan dari yang paling mendesak.
func (s *Service) PesananMenungguPembayaran(ctx context.Context, hariIni time.Time) ([]BarisPesanan, error) {
	rows, err := s.Queries.ListPesanan(ctx, store.ListPesananParams{
		Status: nil,
		Sumber: nil,
		Dari:   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Sampai: time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		return nil, fmt.Errorf("ambil daftar pesanan: %w", err)
	}

	var hasil []BarisPesanan
	for _, row := range rows {
		if row.Status == "batal" || row.Status == "selesai" {
			continue
		}
		if StatusPembayaran(row.TotalHarga, row.TotalBayar) == model.StatusLunas {
			continue
		}
		hasil = append(hasil, s.lengkapiPesanan(row))
	}

	sort.SliceStable(hasil, func(i, j int) bool {
		urut := map[string]int{
			"dp_terlambat":          0,
			"pelunasan_terlambat":   1,
			"dp_jatuh_tempo":        2,
			"pelunasan_jatuh_tempo": 3,
		}
		if urut[hasil[i].LabelJatuhTempo] != urut[hasil[j].LabelJatuhTempo] {
			return urut[hasil[i].LabelJatuhTempo] < urut[hasil[j].LabelJatuhTempo]
		}
		return hasil[i].TanggalAcara.Before(hasil[j].TanggalAcara)
	})
	return hasil, nil
}
