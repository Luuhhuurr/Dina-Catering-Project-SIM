package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"dinacatering/internal/store"
)

// produksiPage menampilkan papan jadwal produksi mingguan — pengganti papan
// tulis di dapur.
func (h *Handler) produksiPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hari := hariIni()
	mulai := hari.AddDate(0, 0, -1)
	selesai := hari.AddDate(0, 0, 13)

	tugas, err := h.svc.Queries.ListTugasRentangTanggal(ctx, store.ListTugasRentangTanggalParams{
		Dari:   mulai,
		Sampai: selesai,
	})
	if err != nil {
		http.Error(w, "gagal memuat jadwal produksi", http.StatusInternalServerError)
		return
	}

	h.render(w, r, "produksi", "Jadwal Produksi", map[string]any{
		"Papan":   kelompokkanPerHari(tugas),
		"HariIni": hari,
		"Mulai":   mulai,
		"Selesai": selesai,
		"Pesan":   r.URL.Query().Get("pesan"),
	})
}

// produksiTable me-render papan jadwal untuk dimuat HTMX.
func (h *Handler) produksiTable(w http.ResponseWriter, r *http.Request) {
	hari := hariIni()
	tugas, err := h.svc.Queries.ListTugasRentangTanggal(r.Context(), store.ListTugasRentangTanggalParams{
		Dari:   hari.AddDate(0, 0, -1),
		Sampai: hari.AddDate(0, 0, 13),
	})
	if err != nil {
		http.Error(w, "gagal memuat jadwal produksi", http.StatusInternalServerError)
		return
	}
	h.renderPartial(w, r, "produksi", "produksi_papan", map[string]any{
		"Papan":   kelompokkanPerHari(tugas),
		"HariIni": hari,
	})
}

// tugasSetPJ menetapkan penanggung jawab sebuah tugas produksi.
func (h *Handler) tugasSetPJ(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id tugas tidak valid", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "permintaan tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.svc.Queries.UpdateTugasPJ(r.Context(), store.UpdateTugasPJParams{
		ID: int32(id),
		Pj: r.FormValue("pj"),
	}); err != nil {
		http.Error(w, "gagal menyimpan penanggung jawab: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/produksi", http.StatusSeeOther)
		return
	}
	toast(w, "Penanggung jawab tugas diperbarui")
	h.produksiTable(w, r)
}

// tugasSetStatus mengubah status pengerjaan tugas produksi.
func (h *Handler) tugasSetStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id tugas tidak valid", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "permintaan tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.svc.Queries.UpdateTugasStatus(r.Context(), store.UpdateTugasStatusParams{
		ID:     int32(id),
		Status: r.FormValue("status"),
	}); err != nil {
		http.Error(w, "gagal mengubah status tugas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/produksi", http.StatusSeeOther)
		return
	}
	toast(w, "Status tugas diperbarui")
	h.produksiTable(w, r)
}
