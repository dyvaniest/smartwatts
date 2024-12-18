package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
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
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	FullName  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// LoginRequest digunakan untuk mengatur input permintaan login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse digunakan untuk mengatur respons dari login
type LoginResponse struct {
	Email   string `json:"email"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type Credential struct {
	Host         string
	Username     string
	Password     string
	DatabaseName string
	Port         int
	Schema       string
}

type BlacklistedToken struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Token     string    `gorm:"unique;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Session struct {
	gorm.Model
	Token  string    `json:"token"`
	Email  string    `json:"email"`
	Expiry time.Time `json:"expiry"`
}
