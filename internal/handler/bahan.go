package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"dinacatering/internal/store"
)

// bahanPage menampilkan daftar bahan baku beserta form tambah dan riwayat harga.
func (h *Handler) bahanPage(w http.ResponseWriter, r *http.Request) {
	daftar, err := h.svc.Queries.ListBahanBaku(r.Context())
	if err != nil {
		http.Error(w, "gagal memuat bahan baku", http.StatusInternalServerError)
		return
	}
	h.render(w, r, "bahan", "Bahan Baku", map[string]any{
		"Daftar": daftar,
		"Pesan":  r.URL.Query().Get("pesan"),
	})
}

// bahanTable me-render tabel bahan baku untuk dimuat HTMX.
func (h *Handler) bahanTable(w http.ResponseWriter, r *http.Request) {
	daftar, err := h.svc.Queries.ListBahanBaku(r.Context())
	if err != nil {
		http.Error(w, "gagal memuat bahan baku", http.StatusInternalServerError)
		return
	}
	h.renderPartial(w, r, "bahan", "bahan_tabel", map[string]any{"Daftar": daftar})
}

// bahanCreate menyimpan bahan baku baru.
func (h *Handler) bahanCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "permintaan tidak valid", http.StatusBadRequest)
		return
	}
	nama := r.FormValue("nama")
	if nama == "" {
		http.Error(w, "nama bahan wajib diisi", http.StatusBadRequest)
		return
	}

	harga, _ := strconv.ParseInt(r.FormValue("harga_terakhir"), 10, 64)
	stok, _ := strconv.ParseFloat(r.FormValue("stok_saat_ini"), 64)

	_, err := h.svc.Queries.CreateBahanBaku(r.Context(), store.CreateBahanBakuParams{
		Nama:          nama,
		Satuan:        orDefault(r.FormValue("satuan"), "kg"),
		StokSaatIni:   stok,
		HargaTerakhir: harga,
	})
	if err != nil {
		http.Error(w, "gagal menyimpan bahan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/bahan", http.StatusSeeOther)
		return
	}
	toast(w, "Bahan "+nama+" tersimpan")
	h.bahanTable(w, r)
}

// bahanDelete menghapus bahan baku.
func (h *Handler) bahanDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id bahan tidak valid", http.StatusBadRequest)
		return
	}
	if err := h.svc.Queries.DeleteBahanBaku(r.Context(), int32(id)); err != nil {
		http.Error(w, "gagal menghapus bahan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/bahan", http.StatusSeeOther)
		return
	}
	toast(w, "Bahan dihapus")
	h.bahanTable(w, r)
}

// hargaBahanRiwayat menampilkan riwayat harga satu bahan beserta form update
// harga. Inilah tombol yang dipakai owner ketika harga di pasar naik.
func (h *Handler) hargaBahanRiwayat(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id bahan tidak valid", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	bahan, err := h.svc.Queries.GetBahanBaku(ctx, int32(id))
	if err != nil {
		http.Error(w, "bahan tidak ditemukan", http.StatusNotFound)
		return
	}
	riwayat, err := h.svc.Queries.ListHargaBahanByBahan(ctx, int32(id))
	if err != nil {
		http.Error(w, "gagal memuat riwayat harga", http.StatusInternalServerError)
		return
	}
	panel := map[string]any{
		"Bahan":   bahan,
		"Riwayat": riwayat,
	}
	if isHTMX(r) {
		h.renderPartial(w, r, "bahan", "harga_riwayat", panel)
		return
	}

	daftar, err := h.svc.Queries.ListBahanBaku(ctx)
	if err != nil {
		http.Error(w, "gagal memuat bahan baku", http.StatusInternalServerError)
		return
	}
	h.render(w, r, "bahan", "Bahan Baku", map[string]any{
		"Daftar": daftar,
		"Pesan":  r.URL.Query().Get("pesan"),
		"Panel":  panel,
	})
}

// hargaBahanTambah menambah riwayat harga baru untuk satu bahan.
func (h *Handler) hargaBahanTambah(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id bahan tidak valid", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "permintaan tidak valid", http.StatusBadRequest)
		return
	}

	harga, err := strconv.ParseInt(r.FormValue("harga_satuan"), 10, 64)
	if err != nil || harga <= 0 {
		http.Error(w, "harga harus berupa angka dan lebih dari nol", http.StatusBadRequest)
		return
	}
	tanggal, err := parseTanggal(r.FormValue("tanggal_berlaku"))
	if err != nil {
		http.Error(w, "tanggal tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.svc.CatatHargaBahan(r.Context(), int32(id), harga, tanggal); err != nil {
		http.Error(w, "gagal menyimpan harga: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/bahan/"+chi.URLParam(r, "id")+"/harga", http.StatusSeeOther)
		return
	}
	toast(w, "Harga baru tersimpan, HPP paket otomatis dihitung ulang")
	h.hargaBahanRiwayat(w, r)
}
