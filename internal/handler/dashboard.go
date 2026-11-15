package handler

import (
	"net/http"
	"time"

	"dinacatering/internal/store"
)

// dashboard menampilkan ringkasan usaha: pesanan berjalan, pembayaran yang
// mendesak, jadwal produksi, dan peringatan margin.
func (h *Handler) dashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hari := hariIni()
	bulanIni := time.Date(hari.Year(), hari.Month(), 1, 0, 0, 0, 0, hari.Location())
	akhirBulan := bulanIni.AddDate(0, 1, -1)

	statusCounts, _ := h.svc.Queries.CountPesananByStatus(ctx)
	sumberCounts, _ := h.svc.Queries.CountPesananBySumber(ctx)
	menunggu, _ := h.svc.PesananMenungguPembayaran(ctx, hari)
	produksi, _ := h.svc.Queries.ListPesananRentangTanggal(ctx, store.ListPesananRentangTanggalParams{
		Dari:   hari,
		Sampai: hari.AddDate(0, 0, 7),
	})
	peringatan, _ := h.svc.PaketDenganMarginDiBawahMinimum(ctx, hari)
	daftarPelanggan, _ := h.svc.Queries.ListPelanggan(ctx, nil)
	labaRugi, _ := h.svc.LaporanLabaRugiPeriode(ctx, bulanIni, akhirBulan)

	menunggu = potong(menunggu, 6)
	produksi = potong(produksi, 6)

	h.render(w, r, "dashboard", "Dashboard", map[string]any{
		"HariIni":           hari,
		"StatusCounts":      statusCounts,
		"SumberCounts":      sumberCounts,
		"MenungguBayar":     menunggu,
		"ProduksiMendatang": produksi,
		"PeringatanMargin":  peringatan,
		"JumlahPelanggan":   len(daftarPelanggan),
		"LabaRugi":          labaRugi,
	})
}

// potong membatasi jumlah baris yang ditampilkan.
func potong[T any](items []T, maksimal int) []T {
	if len(items) > maksimal {
		return items[:maksimal]
	}
	return items
}
