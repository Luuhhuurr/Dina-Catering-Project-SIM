package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"dinacatering/internal/store"
)

// paketPage menampilkan daftar paket lengkap dengan HPP, margin, dan badge
// peringatan margin di bawah 30%.
func (h *Handler) paketPage(w http.ResponseWriter, r *http.Request) {
	hari := hariIni()
	daftar, err := h.svc.DaftarPaket(r.Context(), hari)
	if err != nil {
		http.Error(w, "gagal memuat paket", http.StatusInternalServerError)
		return
	}
	h.render(w, r, "paket", "Paket Menu & HPP", map[string]any{
		"Daftar":  daftar,
		"HariIni": hari,
		"Pesan":   r.URL.Query().Get("pesan"),
	})
}

// paketTable me-render tabel paket untuk dimuat HTMX.
func (h *Handler) paketTable(w http.ResponseWriter, r *http.Request) {
	hari := hariIni()
	daftar, err := h.svc.DaftarPaket(r.Context(), hari)
	if err != nil {
		http.Error(w, "gagal memuat paket", http.StatusInternalServerError)
		return
	}
	h.renderPartial(w, r, "paket", "paket_tabel", map[string]any{
		"Daftar":  daftar,
		"HariIni": hari,
	})
}

// paketCreate menyimpan paket baru.
func (h *Handler) paketCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "permintaan tidak valid", http.StatusBadRequest)
		return
	}
	nama := r.FormValue("nama")
	if nama == "" {
		http.Error(w, "nama paket wajib diisi", http.StatusBadRequest)
		return
	}
	harga, err := strconv.ParseInt(r.FormValue("harga_per_porsi"), 10, 64)
	if err != nil || harga <= 0 {
		http.Error(w, "harga per porsi tidak valid", http.StatusBadRequest)
		return
	}
	tahun, err := strconv.Atoi(r.FormValue("tahun_berlaku"))
	if err != nil {
		tahun = hariIni().Year()
	}

	var hargaWeekend *int64
	if mentah := r.FormValue("harga_per_porsi_weekend"); mentah != "" {
		if nilai, err := strconv.ParseInt(mentah, 10, 64); err == nil && nilai > 0 {
			hargaWeekend = &nilai
		}
	}

	_, err = h.svc.Queries.CreatePaket(r.Context(), store.CreatePaketParams{
		Nama:                 nama,
		HargaPerPorsi:        harga,
		HargaPerPorsiWeekend: hargaWeekend,
		TahunBerlaku:         int32(tahun),
		Deskripsi:            r.FormValue("deskripsi"),
	})
	if err != nil {
		http.Error(w, "gagal menyimpan paket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/paket", http.StatusSeeOther)
		return
	}
	toast(w, "Paket "+nama+" tersimpan")
	h.paketTable(w, r)
}

// paketDelete menghapus paket.
func (h *Handler) paketDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id paket tidak valid", http.StatusBadRequest)
		return
	}
	if err := h.svc.Queries.DeletePaket(r.Context(), int32(id)); err != nil {
		http.Error(w, "gagal menghapus paket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/paket", http.StatusSeeOther)
		return
	}
	toast(w, "Paket dihapus")
	h.paketTable(w, r)
}

// paketIsiTable menampilkan menu-menu yang isi sebuah paket.
func (h *Handler) paketIsiTable(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id paket tidak valid", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	paket, err := h.svc.Queries.GetPaket(ctx, int32(id))
	if err != nil {
		http.Error(w, "paket tidak ditemukan", http.StatusNotFound)
		return
	}
	isi, _ := h.svc.Queries.ListPaketItem(ctx, int32(id))
	menu, _ := h.svc.Queries.ListMenu(ctx)

	panel := map[string]any{
		"Paket": paket,
		"Isi":   isi,
		"Menu":  menu,
	}
	if isHTMX(r) {
		h.renderPartial(w, r, "paket", "paket_isi", panel)
		return
	}

	hari := hariIni()
	daftar, err := h.svc.DaftarPaket(ctx, hari)
	if err != nil {
		http.Error(w, "gagal memuat paket", http.StatusInternalServerError)
		return
	}
	h.render(w, r, "paket", "Paket Menu & HPP", map[string]any{
		"Daftar":  daftar,
		"HariIni": hari,
		"Pesan":   r.URL.Query().Get("pesan"),
		"Panel":   panel,
	})
}

// paketIsiTambah menambahkan menu ke isi paket.
func (h *Handler) paketIsiTambah(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id paket tidak valid", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "permintaan tidak valid", http.StatusBadRequest)
		return
	}
	menuID, err := strconv.ParseInt(r.FormValue("menu_id"), 10, 32)
	if err != nil {
		http.Error(w, "menu wajib dipilih", http.StatusBadRequest)
		return
	}

	_, err = h.svc.Queries.InsertPaketItem(r.Context(), store.InsertPaketItemParams{
		PaketMenuID: int32(id),
		MenuID:      int32(menuID),
	})
	if err != nil {
		http.Error(w, "gagal menambah menu ke paket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/paket/"+chi.URLParam(r, "id")+"/isi", http.StatusSeeOther)
		return
	}
	toast(w, "Menu ditambahkan ke paket")
	h.paketIsiTable(w, r)
}

// paketIsiHapus menghapus menu dari isi paket.
func (h *Handler) paketIsiHapus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id isi paket tidak valid", http.StatusBadRequest)
		return
	}
	if err := h.svc.Queries.DeletePaketItem(r.Context(), int32(id)); err != nil {
		http.Error(w, "gagal menghapus menu dari paket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/paket", http.StatusSeeOther)
		return
	}
	toast(w, "Menu dihapus dari paket")
	h.paketTable(w, r)
}
