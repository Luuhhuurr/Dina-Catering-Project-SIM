package handler

import (
	"net/http"
	"strconv"
	"time"

	"dinacatering/internal/store"
)

// pembelianPage menampilkan catatan belanja bahan di pasar.
func (h *Handler) pembelianPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dari := awalBulan(hariIni())
	sampai := hariIni()

	daftar, err := h.svc.Queries.ListPembelian(ctx, store.ListPembelianParams{Dari: dari, Sampai: sampai})
	if err != nil {
		http.Error(w, "gagal memuat pembelian", http.StatusInternalServerError)
		return
	}
	bahan, _ := h.svc.Queries.ListBahanBaku(ctx)

	h.render(w, r, "pembelian", "Pembelian Bahan", map[string]any{
		"Daftar": daftar,
		"Bahan":  bahan,
		"Dari":   dari,
		"Sampai": sampai,
		"Pesan":  r.URL.Query().Get("pesan"),
	})
}

// pembelianCreate mencatat pembelian bahan di pasar: stok dan harga terkini
// bahan ikut diperbarui.
func (h *Handler) pembelianCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "permintaan tidak valid", http.StatusBadRequest)
		return
	}

	bahanID, err := strconv.ParseInt(r.FormValue("bahan_baku_id"), 10, 32)
	if err != nil {
		http.Error(w, "bahan wajib dipilih", http.StatusBadRequest)
		return
	}
	jumlah, err := strconv.ParseFloat(r.FormValue("jumlah"), 64)
	if err != nil {
		http.Error(w, "jumlah tidak valid", http.StatusBadRequest)
		return
	}
	harga, err := strconv.ParseInt(r.FormValue("harga_satuan"), 10, 64)
	if err != nil || harga <= 0 {
		http.Error(w, "harga satuan tidak valid", http.StatusBadRequest)
		return
	}
	tanggal, err := parseTanggal(r.FormValue("tanggal"))
	if err != nil {
		http.Error(w, "tanggal tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.svc.CatatPembelian(r.Context(), int32(bahanID), jumlah, harga, r.FormValue("pasar"), tanggal); err != nil {
		http.Error(w, "gagal menyimpan pembelian: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/pembelian", http.StatusSeeOther)
		return
	}
	toast(w, "Pembelian tersimpan, stok dan harga bahan diperbarui")

	dari := awalBulan(hariIni())
	sampai := hariIni()
	daftar, err := h.svc.Queries.ListPembelian(r.Context(), store.ListPembelianParams{Dari: dari, Sampai: sampai})
	if err != nil {
		http.Error(w, "gagal memuat pembelian", http.StatusInternalServerError)
		return
	}
	h.renderPartial(w, r, "pembelian", "pembelian_daftar", map[string]any{
		"Daftar": daftar,
		"Dari":   dari,
		"Sampai": sampai,
	})
}

// awalBulan mengembalikan tanggal pertama bulan berjalan.
func awalBulan(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}
