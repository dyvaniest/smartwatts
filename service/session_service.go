package service

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

type SessionService struct {
	store       *sessions.CookieStore
	sessionName string
	cookieName  string
}

// NewSessionService membuat instance baru dari SessionService
func NewSessionService(secretKey string, sessionName string, cookiName string) *SessionService {
	return &SessionService{
		store:       sessions.NewCookieStore([]byte(secretKey)),
		sessionName: sessionName,
		cookieName:  cookiName,
	}
}

// SaveToken menyimpan token ke dalam sesi
func (s *SessionService) SaveToken(w http.ResponseWriter, r *http.Request, token string) error {
	session, err := s.store.Get(r, s.sessionName)
	if err != nil {
		return err
	}

	// Simpan token ke dalam sesi
	session.Values["token"] = token

	if err := session.Save(r, w); err != nil {
		return err
	}

	// Simpan token ke dalam cookie
	http.SetCookie(w, &http.Cookie{
		Name:     s.cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Gunakan true jika HTTPS
	})

	return nil
}

// GetTokenFromSession mengambil token dari sesi
func (s *SessionService) GetTokenFromSession(r *http.Request) (string, error) {
	// Coba ambil token dari sesi
	session, err := s.store.Get(r, s.sessionName)
	if err == nil {
		if token, ok := session.Values["token"].(string); ok && token != "" {
			return token, nil
		}
	}

	// Coba ambil token dari cookie
	cookie, err := r.Cookie(s.cookieName)
	if err == nil {
		return cookie.Value, nil
	}

	// Coba ambil token dari header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}
	}

	return "", errors.New("no token found in session, cookie, or header")
}

// ClearSession menghapus data sesi
func (s *SessionService) ClearToken(w http.ResponseWriter, r *http.Request) error {
	session, err := s.store.Get(r, s.sessionName)
	if err == nil {
		session.Options.MaxAge = -1
		_ = session.Save(r, w)
	}

	// Hapus token dari cookie
	http.SetCookie(w, &http.Cookie{
		Name:     s.cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	return nil
}
