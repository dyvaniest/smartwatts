package main

import (
	"log"
	"net/http"
	"os"

	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/handler"
	"a21hc3NpZ25tZW50/model"
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"a21hc3NpZ25tZW50/service"

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

	// Connect to Database
	db := db.NewDB()
	dbCredential := model.Credential{
		Host:         "localhost",
		Username:     "postgres",
		Password:     "dyvaniest123",
		DatabaseName: "smartwatts",
		Port:         5432,
		Schema:       "public",
	}

	conn, err := db.Connect(&dbCredential)
	if err != nil {
		panic(err)
	}

	conn.AutoMigrate(&model.User{}, &model.Session{})

	userRepo := repository.NewUserRepository(conn)
	sessionRepo := repository.NewSessionRepo(conn)
	userService := service.NewUserService(userRepo)
	sessionService := service.NewSessionService(sessionRepo)

	router := handler.RunServer(tokenAI, userService, sessionService)

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
