package handler

import (
	"a21hc3NpZ25tZW50/service"
	"net/http"

	"github.com/gorilla/mux"
)

func RunServer(tokenAI string, userService *service.UserService, sessionService *service.SessionService) *mux.Router {
	r := mux.NewRouter()

	// User-related endpoints
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		HandleRegister(w, r, userService)
	}).Methods(http.MethodPost)
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		HandleLogin(w, r, userService, sessionService)
	}).Methods(http.MethodPost)
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		HandleLogout(w, r, sessionService)
	}).Methods(http.MethodPost)

	// Data-related endpoints
	r.HandleFunc("/data", HandleData).Methods(http.MethodGet)
	r.HandleFunc("/analytics-energy", HandleAnalyticsEnergy).Methods(http.MethodGet)
	r.HandleFunc("/data", HandleAddData).Methods(http.MethodPost)

	// AI-related endpoints (requires token middleware)
	r.HandleFunc("/analyze-ai", func(w http.ResponseWriter, r *http.Request) {
		HandleAnalyzeAI(w, r, tokenAI)
	}).Methods(http.MethodPost)

	r.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		HandleChat(w, r, tokenAI)
	}).Methods(http.MethodPost)

	return r
}
