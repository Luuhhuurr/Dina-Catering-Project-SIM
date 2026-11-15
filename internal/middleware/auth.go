package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

type contextKey string

const (
	contextKeySession contextKey = "session"
)

// SessionFromContext mengembalikan session yang tersimpan di context request.
func SessionFromContext(ctx context.Context) (Session, bool) {
	session, ok := ctx.Value(contextKeySession).(Session)
	return session, ok
}

// RequireAuth menolak request yang belum login dan mengarahkan ke halaman login.
func RequireAuth(sessions *SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := sessions.Parse(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), contextKeySession, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireGuest mengarahkan user yang sudah login ke dashboard.
func RequireGuest(sessions *SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := sessions.Parse(r); err == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Logging mencatat setiap request beserta waktu dan statusnya.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mulai := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, recorder.status, time.Since(mulai).Round(time.Millisecond))
	})
}

// Recoverer mencegah panic handler menjatuhkan seluruh server.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic pada %s %s: %v", r.Method, r.URL.Path, err)
				http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(status int) {
	s.status = status
	s.ResponseWriter.WriteHeader(status)
}
