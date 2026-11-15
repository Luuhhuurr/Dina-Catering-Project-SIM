// Package handler berisi HTTP handler (chi) untuk seluruh modul aplikasi.
package handler

import (
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"dinacatering/internal/middleware"
	"dinacatering/internal/service"
	"dinacatering/web"
)

// Handler menghubungkan request HTTP dengan service layer.
type Handler struct {
	svc       *service.Service
	sessions  *middleware.SessionManager
	templates map[string]*template.Template
}

// New menyusun Handler beserta template per halaman.
func New(svc *service.Service, sessions *middleware.SessionManager) (*Handler, error) {
	templates, err := buildTemplates()
	if err != nil {
		return nil, err
	}

	return &Handler{svc: svc, sessions: sessions, templates: templates}, nil
}

// buildTemplates mem-parse seluruh halaman. Setiap halaman punya satu set
// template berisi layout dasar, halaman itu sendiri, dan semua partial.
func buildTemplates() (map[string]*template.Template, error) {
	pages := []string{
		"login", "dashboard", "pelanggan", "bahan", "pembelian",
		"menu", "paket", "hpp", "pesanan", "produksi", "laporan",
		"biaya", "pengaturan",
	}

	templates := make(map[string]*template.Template, len(pages))
	for _, page := range pages {
		tpl, err := template.New(page).Funcs(templateFuncs()).ParseFS(web.TemplatesFS,
			"templates/base.layout.html",
			"templates/pages/"+page+".html",
			"templates/partials/*.html",
		)
		if err != nil {
			return nil, err
		}
		templates[page] = tpl
	}
	return templates, nil
}

// Router menyusun seluruh rute aplikasi.
func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logging, middleware.Recoverer)

	r.Method(http.MethodGet, "/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(web.Static()))))

	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return middleware.RequireGuest(h.sessions, next)
		})
		r.Get("/login", h.loginPage)
		r.Post("/login", h.loginSubmit)
		r.Get("/login/", h.loginPage)
		r.Post("/login/", h.loginSubmit)
	})
	r.Post("/logout", h.logout)

	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return middleware.RequireAuth(h.sessions, next)
		})

		r.Get("/", h.dashboard)
		r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		r.Route("/pelanggan", func(r chi.Router) {
			r.Get("/", h.pelangganPage)
			r.Get("/tabel", h.pelangganTable)
			r.Post("/", h.pelangganCreate)
			r.Post("/{id}/hapus", h.pelangganDelete)
		})

		r.Route("/bahan", func(r chi.Router) {
			r.Get("/", h.bahanPage)
			r.Get("/tabel", h.bahanTable)
			r.Post("/", h.bahanCreate)
			r.Post("/{id}/hapus", h.bahanDelete)
			r.Get("/{id}/harga", h.hargaBahanRiwayat)
			r.Post("/{id}/harga", h.hargaBahanTambah)
		})

		r.Route("/pembelian", func(r chi.Router) {
			r.Get("/", h.pembelianPage)
			r.Post("/", h.pembelianCreate)
		})

		r.Route("/menu", func(r chi.Router) {
			r.Get("/", h.menuPage)
			r.Post("/", h.menuCreate)
			r.Post("/{id}/hapus", h.menuDelete)
			r.Get("/{id}/resep", h.resepTable)
			r.Post("/{id}/resep", h.resepCreate)
			r.Post("/resep/{id}/hapus", h.resepDelete)
		})

		r.Route("/paket", func(r chi.Router) {
			r.Get("/", h.paketPage)
			r.Get("/tabel", h.paketTable)
			r.Post("/", h.paketCreate)
			r.Post("/{id}/hapus", h.paketDelete)
			r.Get("/{id}/isi", h.paketIsiTable)
			r.Post("/{id}/isi", h.paketIsiTambah)
			r.Post("/isi/{id}/hapus", h.paketIsiHapus)
		})

		r.Route("/hpp", func(r chi.Router) {
			r.Get("/", h.hppPage)
			r.Get("/hasil", h.hppHasil)
		})

		r.Route("/pesanan", func(r chi.Router) {
			r.Get("/", h.pesananPage)
			r.Get("/tabel", h.pesananTable)
			r.Post("/", h.pesananCreate)
			r.Get("/{id}", h.pesananDetail)
			r.Post("/{id}/status", h.pesananUbahStatus)
			r.Post("/{id}/pembayaran", h.pembayaranCreate)
		})

		r.Route("/produksi", func(r chi.Router) {
			r.Get("/", h.produksiPage)
			r.Get("/tabel", h.produksiTable)
			r.Post("/tugas/{id}/pj", h.tugasSetPJ)
			r.Post("/tugas/{id}/status", h.tugasSetStatus)
		})

		r.Route("/laporan", func(r chi.Router) {
			r.Get("/", h.laporanPage)
		})

		r.Route("/biaya", func(r chi.Router) {
			r.Get("/", h.biayaPage)
			r.Post("/", h.biayaCreate)
		})

		r.Route("/pengaturan", func(r chi.Router) {
			r.Get("/", h.pengaturanPage)
			r.Post("/", h.pengaturanSimpan)
		})
	})

	return r
}

// pageData adalah data yang dikirim ke template dasar.
type pageData struct {
	Judul    string
	Menu     string
	Username string
	Role     string
	Content  any
}

// render merender halaman penuh. Bila request datang dari HTMX, handler
// biasanya memakai renderPartial agar hanya bagian yang berubah yang dikirim.
func (h *Handler) render(w http.ResponseWriter, r *http.Request, halaman string, judul string, data any) {
	tpl, ok := h.templates[halaman]
	if !ok {
		http.Error(w, "halaman tidak ditemukan", http.StatusNotFound)
		return
	}

	username := ""
	role := ""
	if session, ok := middleware.SessionFromContext(r.Context()); ok {
		username = session.Username
		role = session.Role
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(w, "base", pageData{
		Judul:    judul,
		Menu:     halaman,
		Username: username,
		Role:     role,
		Content:  data,
	})
	if err != nil {
		http.Error(w, "gagal merender halaman", http.StatusInternalServerError)
	}
}

// renderPartial merender potongan template untuk dimuat HTMX.
func (h *Handler) renderPartial(w http.ResponseWriter, r *http.Request, halaman string, templateName string, data any) {
	tpl, ok := h.templates[halaman]
	if !ok {
		http.Error(w, "halaman tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(w, templateName, data); err != nil {
		http.Error(w, "gagal merender bagian halaman", http.StatusInternalServerError)
	}
}

// toast mengirim notifikasi singkat ke browser lewat header HX-Trigger.
func toast(w http.ResponseWriter, pesan string) {
	w.Header().Set("HX-Trigger", `{"showToast": "`+pesan+`"}`)
}

// hariIni mengembalikan tanggal hari ini tanpa jam.
func hariIni() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// parseTanggal mengubah string tanggal (YYYY-MM-DD) menjadi time.Time.
func parseTanggal(nilai string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", nilai, time.Local)
}
