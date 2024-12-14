package main

import (
	"log"
	"net/http"
	"os"

	"a21hc3NpZ25tZW50/handler"
	"a21hc3NpZ25tZW50/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Retrieve the Hugging Face token from the environment variables
	tokenAI := os.Getenv("HUGGINGFACE_TOKEN")
	if tokenAI == "" {
		log.Fatal("HUGGINGFACE_TOKEN is not set in the .env file")
	}

	// Get JWT secret key from environment variables
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in the .env file")
	}

	log.Println("JWT Secret loaded successfully!")

	// Set up the router
	router := mux.NewRouter()
	authMiddleware := middleware.AuthenticationMiddleware

	router.HandleFunc("/register", handler.HandleRegister).Methods("POST")
	router.HandleFunc("/login", handler.HandleLogin).Methods("POST")

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(authMiddleware)

	protected.HandleFunc("/data", handler.HandleData).Methods("GET")
	protected.HandleFunc("/analytics-energy", handler.HandleAnalyticsEnergy).Methods("GET", "POST")
	protected.HandleFunc("/analyze-ai", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleAnalyzeAI(w, r, tokenAI)
	}).Methods("POST")
	protected.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleChat(w, r, tokenAI)
	}).Methods("POST")
	protected.HandleFunc("/logout", handler.HandleLogout).Methods("POST")

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow your React app's origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
