package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claim struct {
	Email string `json:"email"`
	jwt.Claims
}

type Inputs struct {
	Table map[string][]string `json:"table"`
	Query string              `json:"query"`
}

type AIRequest struct {
	Inputs Inputs `json:"inputs"`
}

type TapasResponse struct {
	Answer      string   `json:"answer"`
	Coordinates [][]int  `json:"coordinates"`
	Cells       []string `json:"cells"`
	Aggregator  string   `json:"aggregator"`
}

type ChatResponse struct {
	GeneratedText string `json:"generated_text"`
	Answer        string `json:"answer"`
	Status        string `json:"status"`
}

// User struct menyimpan data pengguna
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest digunakan untuk mengatur input permintaan login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse digunakan untuk mengatur respons dari login
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
