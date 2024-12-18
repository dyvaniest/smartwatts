package service

import (
	"a21hc3NpZ25tZW50/model"
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"errors"
	"net/http"
	"strings"
)

type SessionService struct {
	repo repository.SessionsRepository
}

// NewSessionService membuat instance baru dari SessionService
func NewSessionService(repo repository.SessionsRepository) *SessionService {
	return &SessionService{repo: repo}
}

// SaveToken menyimpan token ke dalam sesi di database
func (s *SessionService) SaveToken(email string, token string) error {
	// Membuat session baru
	session := model.Session{
		Email: email,
		Token: token,
	}

	// Simpan session ke repository
	if err := s.repo.AddSessions(session); err != nil {
		return err
	}

	return nil
}

// GetTokenFromSession mengambil token berdasarkan cookie, header, atau repository
func (s *SessionService) GetTokenFromSession(r *http.Request) (string, error) {
	// Ambil token dari cookie
	cookie, err := r.Cookie("auth-token")
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}

	// Ambil token dari header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}
	}

	return "", errors.New("no token found in cookie or header")
}

// ValidateSession memeriksa apakah token masih valid di database
func (s *SessionService) ValidateSession(token string) (model.Session, error) {
	// Cek apakah session tersedia berdasarkan token
	session, err := s.repo.SessionAvailToken(token)
	if err != nil {
		return model.Session{}, err
	}

	return session, nil
}

// ClearToken menghapus token dari repository
func (s *SessionService) ClearToken(token string) error {
	// Hapus session berdasarkan token
	if err := s.repo.DeleteSession(token); err != nil {
		return err
	}

	return nil
}
