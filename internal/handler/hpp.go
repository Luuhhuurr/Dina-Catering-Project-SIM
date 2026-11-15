package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"dinacatering/internal/service"
)

// hppPage menampilkan kalkulator HPP: pemilik memilih paket atau menu dan
// tanggal, lalu sistem menghitung HPP per porsi, margin, dan harga jual
// minimum agar untung minimal 30%.
func (h *Handler) hppPage(w http.ResponseWriter, r *http.Request) {
	tanggal := hariIni()
	if t, err := parseTanggal(r.URL.Query().Get("tanggal")); err == nil {
		tanggal = t
	}
	h.render(w, r, "hpp", "Kalkulator HPP", h.dataHPP(r.Context(), tanggal, nil, "", 0, 0))
}

// dataHPP menyusun data halaman kalkulator HPP. Key "Panel" berisi hasil
// perhitungan beserta jenisnya ("paket" atau "menu") sehingga halaman penuh
// dan partial HTMX memakai sumber data yang sama.
func (h *Handler) dataHPP(ctx context.Context, tanggal time.Time, panel map[string]any, jenis string, paketTerpilih, menuTerpilih int32) map[string]any {
	paket, _ := h.svc.Queries.ListPaket(ctx)
	menu, _ := h.svc.Queries.ListMenu(ctx)
	persen, _ := h.svc.PersenBiayaOperasional(ctx)

	data := map[string]any{
		"Paket":         paket,
		"Menu":          menu,
		"Persen":        persen,
		"HariIni":       hariIni(),
		"Tanggal":       tanggal.Format("2006-01-02"),
		"Jenis":         jenis,
		"PaketTerpilih": paketTerpilih,
		"MenuTerpilih":  menuTerpilih,
	}
	if panel != nil {
		data["Panel"] = panel
	}
	return data
}

// hppHasil menghitung dan menampilkan hasil kalkulator HPP. Permintaan HTMX
// menerima potongan HTML; permintaan biasa (mis. JavaScript dimatikan)
// menerima halaman penuh berisi hasil yang sama.
func (h *Handler) hppHasil(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tanggal := hariIni()
	if t, err := parseTanggal(r.URL.Query().Get("tanggal")); err == nil {
		tanggal = t
	}
	persen, _ := h.svc.PersenBiayaOperasional(ctx)

	if paketID := r.URL.Query().Get("paket_id"); paketID != "" {
		id, err := strconv.ParseInt(paketID, 10, 32)
		if err != nil {
			http.Error(w, "id paket tidak valid", http.StatusBadRequest)
			return
		}
		paket, err := h.svc.Queries.GetPaket(ctx, int32(id))
		if err != nil {
			http.Error(w, "paket tidak ditemukan", http.StatusNotFound)
			return
		}
		hpp, err := h.svc.HitungHPPPaket(ctx, h.svc.Queries, paket.ID, tanggal, persen)
		if err != nil {
			http.Error(w, "gagal menghitung HPP: "+err.Error(), http.StatusInternalServerError)
			return
		}

		hargaJual := service.HargaPerPorsiEfektif(paket.HargaPerPorsi, paket.HargaPerPorsiWeekend, tanggal)
		margin := service.MarginPaket(hargaJual, hpp.HPPPerPorsi)

		panel := map[string]any{
			"Paket":        paket,
			"HPP":          hpp,
			"HargaJual":    hargaJual,
			"Margin":       margin,
			"HargaMinimum": service.HargaJualMinimum30Persen(hpp.HPPPerPorsi),
			"Peringatan":   service.MarginDiBawahMinimum(margin),
			"Tanggal":      tanggal,
		}
		if isHTMX(r) {
			h.renderPartial(w, r, "hpp", "hpp_hasil_paket", panel)
			return
		}
		h.render(w, r, "hpp", "Kalkulator HPP", h.dataHPP(ctx, tanggal, panel, "paket", paket.ID, 0))
		return
	}

	menuID := r.URL.Query().Get("menu_id")
	if menuID == "" {
		if isHTMX(r) {
			h.renderPartial(w, r, "hpp", "hpp_hasil_paket", map[string]any{"Kosong": true})
			return
		}
		h.render(w, r, "hpp", "Kalkulator HPP", h.dataHPP(ctx, tanggal, nil, "", 0, 0))
		return
	}
	id, err := strconv.ParseInt(menuID, 10, 32)
	if err != nil {
		http.Error(w, "id menu tidak valid", http.StatusBadRequest)
		return
	}
	menu, err := h.svc.Queries.GetMenu(ctx, int32(id))
	if err != nil {
		http.Error(w, "menu tidak ditemukan", http.StatusNotFound)
		return
	}
	hpp, err := h.svc.HitungHPPMenu(ctx, h.svc.Queries, int32(id), tanggal, persen)
	if err != nil {
		http.Error(w, "gagal menghitung HPP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	panel := map[string]any{
		"HPP":          hpp,
		"Tanggal":      tanggal,
		"HargaMinimum": service.HargaJualMinimum30Persen(hpp.HPPPerPorsi),
	}
	if isHTMX(r) {
		h.renderPartial(w, r, "hpp", "hpp_hasil_menu", panel)
		return
	}
	h.render(w, r, "hpp", "Kalkulator HPP", h.dataHPP(ctx, tanggal, panel, "menu", 0, menu.ID))
}
