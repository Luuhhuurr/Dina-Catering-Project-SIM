// Package middleware berisi middleware HTTP: logging, recover, dan autentikasi.
package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	sessionCookieName = "dina_session"
	sessionTTL        = 12 * time.Hour
)

// Session adalah session bertanda tangan yang disimpan di dalam cookie.
// Isinya hanya user ID dan waktu kedaluwarsa.
type Session struct {
	UserID    int32
	Username  string
	Role      string
	ExpiresAt time.Time
}

// SessionManager membuat dan memverifikasi cookie session.
type SessionManager struct {
	secret []byte
}

// NewSessionManager menyusun SessionManager dari SESSION_SECRET.
func NewSessionManager(secret string) *SessionManager {
	return &SessionManager{secret: []byte(secret)}
}

var errSessionInvalid = errors.New("session tidak valid")

// Issue membuat cookie session untuk user tertentu.
func (m *SessionManager) Issue(w http.ResponseWriter, userID int32, username, role string) {
	now := time.Now()
	payload := strings.Join([]string{
		strconv.FormatInt(int64(userID), 10),
		base64.RawURLEncoding.EncodeToString([]byte(username)),
		base64.RawURLEncoding.EncodeToString([]byte(role)),
		strconv.FormatInt(now.Add(sessionTTL).Unix(), 10),
	}, "|")

	encoded := base64.RawURLEncoding.EncodeToString([]byte(payload))
	signature := m.sign(encoded)

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    encoded + "." + signature,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  now.Add(sessionTTL),
	})
}

// Clear menghapus cookie session.
func (m *SessionManager) Clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

// Parse memverifikasi cookie session dan mengembalikan isi session.
func (m *SessionManager) Parse(r *http.Request) (Session, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return Session{}, errSessionInvalid
	}

	encoded, signature, ok := strings.Cut(cookie.Value, ".")
	if !ok {
		return Session{}, errSessionInvalid
	}
	if !hmac.Equal([]byte(signature), []byte(m.sign(encoded))) {
		return Session{}, errSessionInvalid
	}

	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return Session{}, errSessionInvalid
	}
	parts := strings.Split(string(decoded), "|")
	if len(parts) != 4 {
		return Session{}, errSessionInvalid
	}

	userID, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return Session{}, errSessionInvalid
	}
	username, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Session{}, errSessionInvalid
	}
	role, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return Session{}, errSessionInvalid
	}
	expUnix, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return Session{}, errSessionInvalid
	}

	expired := time.Now().After(time.Unix(expUnix, 0))
	if expired {
		return Session{}, errSessionInvalid
	}

	return Session{
		UserID:    int32(userID),
		Username:  string(username),
		Role:      string(role),
		ExpiresAt: time.Unix(expUnix, 0),
	}, nil
}

func (m *SessionManager) sign(payload string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
