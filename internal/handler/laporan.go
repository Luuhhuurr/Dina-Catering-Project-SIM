package handler

import (
	"net/http"
	"strconv"
	"time"

	"dinacatering/internal/store"
)

// laporanPage menampilkan laba rugi, rekap pesanan, dan margin per pesanan.
func (h *Handler) laporanPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hari := hariIni()

	dari := awalBulan(hari)
	sampai := hari
	if nilai := r.URL.Query().Get("bulan"); nilai != "" {
		if t, err := parseTanggal(nilai + "-01"); err == nil {
			dari = t
			sampai = t.AddDate(0, 1, -1)
		}
	}

	labaRugi, err := h.svc.LaporanLabaRugiPeriode(ctx, dari, sampai)
	if err != nil {
		http.Error(w, "gagal memuat laporan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	statusCounts, _ := h.svc.Queries.CountPesananByStatus(ctx)
	sumberCounts, _ := h.svc.Queries.CountPesananBySumber(ctx)

	pesanan, _ := h.svc.Queries.ListPesanan(ctx, store.ListPesananParams{
		Status: nil,
		Sumber: nil,
		Dari:   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Sampai: time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	margin, _ := h.svc.MarginSeluruhPesanan(ctx, pesanan, hari)

	peringatanMargin, _ := h.svc.PaketDenganMarginDiBawahMinimum(ctx, hari)
	bahan, _ := h.svc.Queries.ListBahanBaku(ctx)

	h.render(w, r, "laporan", "Laporan", map[string]any{
		"LabaRugi":         labaRugi,
		"Dari":             dari,
		"Sampai":           sampai,
		"Bulan":            hari.Format("2006-01"),
		"StatusCounts":     statusCounts,
		"SumberCounts":     sumberCounts,
		"Margin":           margin,
		"PeringatanMargin": peringatanMargin,
		"Bahan":            bahan,
	})
}

// biayaPage menampilkan catatan biaya operasional.
func (h *Handler) biayaPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dari := awalBulan(hariIni())
	sampai := hariIni()

	daftar, err := h.svc.Queries.ListBiaya(ctx, store.ListBiayaParams{Dari: dari, Sampai: sampai})
	if err != nil {
		http.Error(w, "gagal memuat biaya", http.StatusInternalServerError)
		return
	}
	total, _ := h.svc.Queries.TotalBiayaRentang(ctx, store.TotalBiayaRentangParams{Dari: dari, Sampai: sampai})

	h.render(w, r, "biaya", "Biaya Operasional", map[string]any{
		"Daftar": daftar,
		"Total":  total,
		"Dari":   dari,
		"Sampai": sampai,
	})
}

// biayaCreate menyimpan biaya operasional baru.
func (h *Handler) biayaCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "permintaan tidak valid", http.StatusBadRequest)
		return
	}
	jumlah, err := strconv.ParseInt(r.FormValue("jumlah"), 10, 64)
	if err != nil || jumlah <= 0 {
		http.Error(w, "nominal biaya tidak valid", http.StatusBadRequest)
		return
	}
	tanggal, err := parseTanggal(r.FormValue("tanggal"))
	if err != nil {
		tanggal = hariIni()
	}

	_, err = h.svc.Queries.CreateBiaya(r.Context(), store.CreateBiayaParams{
		Kategori:  r.FormValue("kategori"),
		Deskripsi: r.FormValue("deskripsi"),
		Jumlah:    jumlah,
		Tanggal:   tanggal,
	})
	if err != nil {
		http.Error(w, "gagal menyimpan biaya: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/biaya", http.StatusSeeOther)
		return
	}
	toast(w, "Biaya operasional tersimpan")

	dari := awalBulan(hariIni())
	sampai := hariIni()
	daftar, err := h.svc.Queries.ListBiaya(r.Context(), store.ListBiayaParams{Dari: dari, Sampai: sampai})
	if err != nil {
		http.Error(w, "gagal memuat biaya", http.StatusInternalServerError)
		return
	}
	total, _ := h.svc.Queries.TotalBiayaRentang(r.Context(), store.TotalBiayaRentangParams{Dari: dari, Sampai: sampai})
	h.renderPartial(w, r, "biaya", "biaya_daftar", map[string]any{
		"Daftar": daftar,
		"Total":  total,
	})
}

// pengaturanPage menampilkan pengaturan perhitungan (persentase biaya
// operasional yang ditambahkan ke HPP bahan).
func (h *Handler) pengaturanPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	persen, err := h.svc.PersenBiayaOperasional(ctx)
	if err != nil {
		http.Error(w, "gagal memuat pengaturan", http.StatusInternalServerError)
		return
	}
	daftar, _ := h.svc.Queries.ListPengaturan(ctx)

	h.render(w, r, "pengaturan", "Pengaturan", map[string]any{
		"Persen": persen,
		"Daftar": daftar,
		"Pesan":  r.URL.Query().Get("pesan"),
	})
}

// pengaturanSimpan menyimpan persentase biaya operasional.
func (h *Handler) pengaturanSimpan(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "permintaan tidak valid", http.StatusBadRequest)
		return
	}
	persen, err := strconv.ParseFloat(r.FormValue("biaya_operasional_persen"), 64)
	if err != nil || persen < 0 || persen > 100 {
		http.Error(w, "persentase harus antara 0 dan 100", http.StatusBadRequest)
		return
	}

	if err := h.svc.SetPersenBiayaOperasional(r.Context(), persen); err != nil {
		http.Error(w, "gagal menyimpan pengaturan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/pengaturan", http.StatusSeeOther)
		return
	}
	toast(w, "Pengaturan tersimpan, HPP dihitung ulang")
	h.pengaturanPartial(w, r)
}

// pengaturanPartial merender ulang formulir pengaturan.
func (h *Handler) pengaturanPartial(w http.ResponseWriter, r *http.Request) {
	persen, err := h.svc.PersenBiayaOperasional(r.Context())
	if err != nil {
		http.Error(w, "gagal memuat pengaturan", http.StatusInternalServerError)
		return
	}
	daftar, _ := h.svc.Queries.ListPengaturan(r.Context())
	h.renderPartial(w, r, "pengaturan", "pengaturan_form", map[string]any{
		"Persen": persen,
		"Daftar": daftar,
	})
}
