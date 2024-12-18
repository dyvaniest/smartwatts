package handler

import (
	"encoding/json"
	"net/http"

	"a21hc3NpZ25tZW50/model"
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"a21hc3NpZ25tZW50/service"

	"gorm.io/gorm"
)

var userRepo = repository.NewUserRepository(&gorm.DB{})
var userService = service.NewUserService(userRepo)
var sessionRepo = repository.NewSessionRepo(&gorm.DB{})
var sessionService = service.NewSessionService(sessionRepo)

func HandleRegister(w http.ResponseWriter, r *http.Request, userService *service.UserService) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := userService.RegisterUser(user)
	if err != nil {
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}

func HandleLogin(w http.ResponseWriter, r *http.Request, userService *service.UserService, sessionService *service.SessionService) {
	var loginRequest model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	loginResponse, err := userService.LoginUser(loginRequest)
	if err != nil {
		http.Error(w, "Failed to login: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Simpan token ke dalam sesi
	err = sessionService.SaveToken(loginResponse.Email, loginResponse.Token)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	// Simpan token ke dalam cookie untuk frontend
	http.SetCookie(w, &http.Cookie{
		Name:     "auth-token",
		Value:    loginResponse.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	json.NewEncoder(w).Encode(loginResponse)
}

func HandleLogout(w http.ResponseWriter, r *http.Request, sessionService *service.SessionService) {
	// Ambil token dari session atau cookie
	token, err := sessionService.GetTokenFromSession(r)
	if err != nil {
		http.Error(w, "Failed to retrieve token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	err = sessionService.ClearToken(token)
	if err != nil {
		http.Error(w, "Failed to clear session", http.StatusInternalServerError)
		return
	}

	// Hapus cookie dari browser
	http.SetCookie(w, &http.Cookie{
		Name:     "auth-token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1, // Expire immediately
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully logged out",
	})
}
