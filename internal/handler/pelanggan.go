package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"dinacatering/internal/store"
)

// pelangganPage menampilkan daftar pelanggan beserta form tambah.
func (h *Handler) pelangganPage(w http.ResponseWriter, r *http.Request) {
	cari := r.URL.Query().Get("cari")
	daftar, err := h.svc.Queries.ListPelanggan(r.Context(), &cari)
	if err != nil {
		http.Error(w, "gagal memuat pelanggan", http.StatusInternalServerError)
		return
	}
	h.render(w, r, "pelanggan", "Pelanggan", map[string]any{
		"Daftar": daftar,
		"Cari":   cari,
		"Pesan":  r.URL.Query().Get("pesan"),
	})
}

// pelangganTable me-render tabel pelanggan untuk dimuat HTMX.
func (h *Handler) pelangganTable(w http.ResponseWriter, r *http.Request) {
	cari := r.URL.Query().Get("cari")
	daftar, err := h.svc.Queries.ListPelanggan(r.Context(), &cari)
	if err != nil {
		http.Error(w, "gagal memuat pelanggan", http.StatusInternalServerError)
		return
	}
	h.renderPartial(w, r, "pelanggan", "pelanggan_tabel", map[string]any{
		"Daftar": daftar,
		"Cari":   cari,
	})
}

// pelangganCreate menyimpan pelanggan baru.
func (h *Handler) pelangganCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "permintaan tidak valid", http.StatusBadRequest)
		return
	}

	nama := r.FormValue("nama")
	if nama == "" {
		http.Error(w, "nama pelanggan wajib diisi", http.StatusBadRequest)
		return
	}

	_, err := h.svc.Queries.CreatePelanggan(r.Context(), store.CreatePelangganParams{
		Nama:    nama,
		Kontak:  r.FormValue("kontak"),
		Alamat:  r.FormValue("alamat"),
		Catatan: r.FormValue("catatan"),
	})
	if err != nil {
		http.Error(w, "gagal menyimpan pelanggan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/pelanggan", http.StatusSeeOther)
		return
	}
	toast(w, "Pelanggan "+nama+" tersimpan")
	h.pelangganTable(w, r)
}

// pelangganDelete menghapus pelanggan.
func (h *Handler) pelangganDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id pelanggan tidak valid", http.StatusBadRequest)
		return
	}
	if err := h.svc.Queries.DeletePelanggan(r.Context(), int32(id)); err != nil {
		http.Error(w, "gagal menghapus pelanggan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/pelanggan", http.StatusSeeOther)
		return
	}
	toast(w, "Pelanggan dihapus")
	h.pelangganTable(w, r)
}

// isHTMX menandai request yang berasal dari HTMX.
func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") != ""
}
