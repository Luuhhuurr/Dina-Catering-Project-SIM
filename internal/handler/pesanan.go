package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"dinacatering/internal/model"
	"dinacatering/internal/service"
	"dinacatering/internal/store"
)

// pesananPage menampilkan daftar pesanan beserta form pemesanan baru.
func (h *Handler) pesananPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hari := hariIni()
	status := r.URL.Query().Get("status")
	sumber := r.URL.Query().Get("sumber")

	daftar, err := h.svc.DaftarPesanan(ctx, status, sumber, awalBulan(hari), hari.AddDate(1, 0, 0))
	if err != nil {
		http.Error(w, "gagal memuat pesanan", http.StatusInternalServerError)
		return
	}
	menunggu, _ := h.svc.PesananMenungguPembayaran(ctx, hari)

	h.render(w, r, "pesanan", "Pesanan", map[string]any{
		"Daftar":        daftar,
		"MenungguBayar": menunggu,
		"Status":        status,
		"Sumber":        sumber,
		"Pelanggan":     h.daftarPelanggan(ctx),
		"Paket":         h.daftarPaket(ctx),
		"StatusPesanan": model.StatusPesananList,
		"SumberPesanan": model.SumberPesananList,
		"HariIni":       hari,
		"Pesan":         r.URL.Query().Get("pesan"),
	})
}

// pesananTable me-render tabel pesanan untuk dimuat HTMX.
func (h *Handler) pesananTable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hari := hariIni()
	status := r.URL.Query().Get("status")
	sumber := r.URL.Query().Get("sumber")

	daftar, err := h.svc.DaftarPesanan(ctx, status, sumber, awalBulan(hari), hari.AddDate(1, 0, 0))
	if err != nil {
		http.Error(w, "gagal memuat pesanan", http.StatusInternalServerError)
		return
	}
	h.renderPartial(w, r, "pesanan", "pesanan_tabel", map[string]any{
		"Daftar":        daftar,
		"Status":        status,
		"Sumber":        sumber,
		"StatusPesanan": model.StatusPesananList,
		"SumberPesanan": model.SumberPesananList,
	})
}

// pesananCreate membuat pesanan baru: total harga dan tenggat DP serta
// pelunasan dihitung otomatis.
func (h *Handler) pesananCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "permintaan tidak valid", http.StatusBadRequest)
		return
	}

	pelangganID, err := strconv.ParseInt(r.FormValue("pelanggan_id"), 10, 32)
	if err != nil {
		http.Error(w, "pelanggan wajib dipilih", http.StatusBadRequest)
		return
	}
	jumlahTamu, err := strconv.Atoi(r.FormValue("jumlah_tamu"))
	if err != nil || jumlahTamu <= 0 {
		http.Error(w, "jumlah tamu harus lebih dari nol", http.StatusBadRequest)
		return
	}
	tanggalAcara, err := parseTanggal(r.FormValue("tanggal_acara"))
	if err != nil {
		http.Error(w, "tanggal acara tidak valid", http.StatusBadRequest)
		return
	}

	masakDiLokasi := r.FormValue("masak_di_lokasi") == "on"

	pesananID, err := h.svc.BuatPesanan(r.Context(), service.PesananInput{
		PelangganID:   int32(pelangganID),
		NamaPaket:     r.FormValue("nama_paket"),
		JumlahTamu:    jumlahTamu,
		TanggalAcara:  tanggalAcara,
		Lokasi:        r.FormValue("lokasi"),
		SumberPesanan: r.FormValue("sumber_pesanan"),
		Catatan:       r.FormValue("catatan"),
		MasakDiLokasi: masakDiLokasi,
	})
	if err != nil {
		http.Error(w, "gagal menyimpan pesanan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/pesanan/"+strconv.Itoa(int(pesananID)), http.StatusSeeOther)
		return
	}
	toast(w, "Pesanan tersimpan, tenggat DP H-30 dan pelunasan H-7 dihitung otomatis")
	h.pesananTable(w, r)
}

// pesananDetail menampilkan rincian satu pesanan: pembayaran, tenggat, dan
// margin untungnya.
func (h *Handler) pesananDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id pesanan tidak valid", http.StatusBadRequest)
		return
	}
	data, err := h.dataPesananDetail(r, int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.render(w, r, "pesanan", "Detail Pesanan", data)
}

// dataPesananDetail mengumpulkan seluruh data yang dipakai halaman detail
// pesanan.
func (h *Handler) dataPesananDetail(r *http.Request, id int32) (map[string]any, error) {
	ctx := r.Context()
	pesanan, err := h.svc.Queries.GetPesanan(ctx, id)
	if err != nil {
		return nil, errors.New("pesanan tidak ditemukan")
	}
	pembayaran, _ := h.svc.Queries.ListPembayaranByPesanan(ctx, id)
	tugas, _ := h.svc.Queries.ListTugasByPesanan(ctx, id)
	margin, _ := h.svc.MarginPerPesanan(ctx, id)

	dp, pelunasan := service.JadwalPembayaran(pesanan.TanggalAcara)

	return map[string]any{
		"Pesanan":          pesanan,
		"Pembayaran":       pembayaran,
		"Tugas":            tugas,
		"Margin":           margin,
		"TenggatDP":        dp,
		"TenggatLunas":     pelunasan,
		"StatusPembayaran": service.StatusPembayaran(pesanan.TotalHarga, pesanan.TotalBayar),
		"StatusPesanan":    model.StatusPesananList,
		"Metode":           model.MetodePembayaranList,
		"NilaiDP":          service.DPAmount(pesanan.TotalHarga),
		"NilaiPelunasan":   pesanan.TotalHarga - service.DPAmount(pesanan.TotalHarga),
		"Detail":           true,
		"HariIni":          hariIni(),
	}, nil
}

// pesananUbahStatus mengubah status alur pesanan.
func (h *Handler) pesananUbahStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id pesanan tidak valid", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "permintaan tidak valid", http.StatusBadRequest)
		return
	}
	statusBaru := r.FormValue("status")
	if statusBaru == "" {
		http.Error(w, "status wajib diisi", http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateStatusPesanan(r.Context(), int32(id), statusBaru); err != nil {
		http.Error(w, "gagal mengubah status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/pesanan", http.StatusSeeOther)
		return
	}
	toast(w, "Status pesanan diubah menjadi "+model.LabelStatusPesanan[statusBaru])

	if r.FormValue("dari") == "detail" {
		data, err := h.dataPesananDetail(r, int32(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		h.renderPartial(w, r, "pesanan", "pesanan_detail", data)
		return
	}
	h.pesananTable(w, r)
}

// pembayaranCreate mencatat pembayaran DP atau pelunasan.
func (h *Handler) pembayaranCreate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "id pesanan tidak valid", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "permintaan tidak valid", http.StatusBadRequest)
		return
	}
	jumlah, err := strconv.ParseInt(r.FormValue("jumlah"), 10, 64)
	if err != nil || jumlah <= 0 {
		http.Error(w, "nominal pembayaran tidak valid", http.StatusBadRequest)
		return
	}
	tanggal, err := parseTanggal(r.FormValue("tanggal"))
	if err != nil {
		tanggal = hariIni()
	}

	tipe := r.FormValue("tipe")
	if tipe != model.TipePelunasan {
		tipe = model.TipeDP
	}
	metode := orDefault(r.FormValue("metode"), "transfer")

	if err := h.svc.CatatPembayaran(r.Context(), int32(id), jumlah, tipe, metode, tanggal); err != nil {
		http.Error(w, "gagal menyimpan pembayaran: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !isHTMX(r) {
		http.Redirect(w, r, "/pesanan/"+strconv.Itoa(int(id)), http.StatusSeeOther)
		return
	}
	toast(w, "Pembayaran "+tipe+" tercatat")
	data, err := h.dataPesananDetail(r, int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.renderPartial(w, r, "pesanan", "pesanan_detail", data)
}

// daftarPelanggan mengembalikan daftar pelanggan untuk form pesanan.
func (h *Handler) daftarPelanggan(ctx context.Context) []store.ListPelangganRow {
	daftar, _ := h.svc.Queries.ListPelanggan(ctx, nil)
	return daftar
}

// daftarPaket mengembalikan daftar paket untuk form pesanan.
func (h *Handler) daftarPaket(ctx context.Context) []store.ListPaketRow {
	daftar, _ := h.svc.Queries.ListPaket(ctx)
	return daftar
}
