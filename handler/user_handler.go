package handler

import (
	"encoding/json"
	"net/http"

	"a21hc3NpZ25tZW50/model"
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"a21hc3NpZ25tZW50/service"
)

var userRepo = repository.NewUserRepository()
var userService = service.NewUserService(userRepo)
var sessionService = service.NewSessionService("my-secret-key", "auth-session", "auth-token")

func HandleRegister(w http.ResponseWriter, r *http.Request) {
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

// Handle user login
func HandleLogin(w http.ResponseWriter, r *http.Request) {
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
	err = sessionService.SaveToken(w, r, loginResponse.Token)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(loginResponse)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	err := sessionService.ClearToken(w, r)
	if err != nil {
		http.Error(w, "Failed to clear session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully logged out",
	})
}
