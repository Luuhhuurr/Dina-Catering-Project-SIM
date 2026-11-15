package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"dinacatering/internal/store"
)

// menuPage menampilkan daftar menu dan panel resep.
func (h *Handler) menuPage(w http.ResponseWriter, r *http.Request) {
	daftar, err := h.svc.Queries.ListMenu(r.Context())
	if err != nil {
		http.Error(w, "gagal memuat menu", http.StatusInternalServerError)
		return
	}
	h.render(w, r, "menu", "Menu & Resep", map[string]any{
		"Daftar": daftar,
		"Pesan":  r.URL.Query().Get("pesan"),
	})
}

// menuCreate menyimpan menu baru.
func (h *Handler) menuCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "permintaan tidak valid", http.StatusBadRequest)
		return
	}
	nama := r.FormValue("nama")
	if nama == "" {
		http.Error(w, "nama menu wajib diisi", http.StatusBadRequest)
		return
	}

	_, err := h.svc.Queries.CreateMenu(r.Context(), store.CreateMenuParams{
		Nama:      nama,
		Deskripsi: r.FormValue("deskripsi"),
	})
	if err != nil {
		http.Error(w, "gagal menyimpan menu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/menu", http.StatusSeeOther)
		return
	}
	toast(w, "Menu "+nama+" tersimpan")
	h.menuTable(w, r)
}

// menuTable me-render daftar menu untuk dimuat HTMX.
func (h *Handler) menuTable(w http.ResponseWriter, r *http.Request) {
	daftar, err := h.svc.Queries.ListMenu(r.Context())
	if err != nil {
		http.Error(w, "gagal memuat menu", http.StatusInternalServerError)
		return
	}
	h.renderPartial(w, r, "menu", "menu_tabel", map[string]any{"Daftar": daftar})
}

// menuDelete menghapus menu beserta resepnya.
func (h *Handler) menuDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id menu tidak valid", http.StatusBadRequest)
		return
	}
	if err := h.svc.Queries.DeleteMenu(r.Context(), int32(id)); err != nil {
		http.Error(w, "gagal menghapus menu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/menu", http.StatusSeeOther)
		return
	}
	toast(w, "Menu dihapus")
	h.menuTable(w, r)
}

// resepTable menampilkan resep satu menu beserta HPP per porsinya.
func (h *Handler) resepTable(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id menu tidak valid", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	menu, err := h.svc.Queries.GetMenu(ctx, int32(id))
	if err != nil {
		http.Error(w, "menu tidak ditemukan", http.StatusNotFound)
		return
	}
	resep, _ := h.svc.Queries.ListResepByMenu(ctx, int32(id))
	bahan, _ := h.svc.Queries.ListBahanBaku(ctx)
	persen, _ := h.svc.PersenBiayaOperasional(ctx)
	hpp, _ := h.svc.HitungHPPMenu(ctx, h.svc.Queries, int32(id), hariIni(), persen)

	panel := map[string]any{
		"Menu":   menu,
		"Resep":  resep,
		"Bahan":  bahan,
		"HPP":    hpp,
		"Persen": persen,
	}
	if isHTMX(r) {
		h.renderPartial(w, r, "menu", "resep_panel", panel)
		return
	}

	daftar, err := h.svc.Queries.ListMenu(ctx)
	if err != nil {
		http.Error(w, "gagal memuat menu", http.StatusInternalServerError)
		return
	}
	h.render(w, r, "menu", "Menu & Resep", map[string]any{
		"Daftar": daftar,
		"Pesan":  r.URL.Query().Get("pesan"),
		"Panel":  panel,
	})
}

// resepCreate menambahkan satu bahan ke resep menu.
func (h *Handler) resepCreate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id menu tidak valid", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "permintaan tidak valid", http.StatusBadRequest)
		return
	}

	bahanID, err := strconv.ParseInt(r.FormValue("bahan_baku_id"), 10, 32)
	if err != nil {
		http.Error(w, "bahan wajib dipilih", http.StatusBadRequest)
		return
	}
	jumlah, err := strconv.ParseFloat(r.FormValue("jumlah_per_porsi"), 64)
	if err != nil || jumlah <= 0 {
		http.Error(w, "takaran harus berupa angka dan lebih dari nol", http.StatusBadRequest)
		return
	}

	_, err = h.svc.Queries.InsertResepItem(r.Context(), store.InsertResepItemParams{
		MenuID:         int32(id),
		BahanBakuID:    int32(bahanID),
		JumlahPerPorsi: jumlah,
		Satuan:         r.FormValue("satuan"),
	})
	if err != nil {
		http.Error(w, "gagal menyimpan resep: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/menu/"+chi.URLParam(r, "id")+"/resep", http.StatusSeeOther)
		return
	}
	toast(w, "Resep diperbarui, HPP dihitung ulang")
	h.resepTable(w, r)
}

// resepDelete menghapus satu bahan dari resep menu.
func (h *Handler) resepDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id resep tidak valid", http.StatusBadRequest)
		return
	}
	if err := h.svc.Queries.DeleteResepItem(r.Context(), int32(id)); err != nil {
		http.Error(w, "gagal menghapus resep: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/menu", http.StatusSeeOther)
		return
	}
	toast(w, "Resep diperbarui, HPP dihitung ulang")
	h.menuTable(w, r)
}
