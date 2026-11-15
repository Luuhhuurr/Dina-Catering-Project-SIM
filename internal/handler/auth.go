package handler

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// loginPage menampilkan halaman login.
func (h *Handler) loginPage(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "login", "Masuk", map[string]any{"Error": ""})
}

// loginSubmit memproses login owner/admin.
func (h *Handler) loginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.gagalLogin(w, r, "Permintaan tidak valid.")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.svc.FindUserByUsername(r.Context(), username)
	if err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	}
	if err != nil {
		h.gagalLogin(w, r, "Username atau password salah.")
		return
	}

	h.sessions.Issue(w, user.ID, user.Username, user.Role)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// gagalLogin menampilkan kembali kotak login dengan pesan kesalahan.
func (h *Handler) gagalLogin(w http.ResponseWriter, r *http.Request, pesan string) {
	data := map[string]any{"Error": pesan}
	if isHTMX(r) {
		h.renderPartial(w, r, "login", "login_kotak", data)
		return
	}
	h.render(w, r, "login", "Masuk", data)
}

// logout menghapus session dan mengembalikan ke halaman login.
func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	h.sessions.Clear(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
