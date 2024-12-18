package repository

import (
	"a21hc3NpZ25tZW50/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// AddUser menambahkan pengguna baru ke database
func (r *UserRepository) AddUser(user model.User) error {
	var existingUser model.User
	err := r.db.Where("email = ?", user.Email).First(&existingUser).Error
	if err == nil {
		return fmt.Errorf("email already registered")
	}

	// Tambahkan pengguna ke database
	user.CreatedAt = time.Now()
	if err := r.db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByEmail mengambil pengguna berdasarkan email
func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// BlacklistToken menambahkan token ke tabel blacklist
func (r *UserRepository) BlacklistToken(token string) error {
	blacklist := model.BlacklistedToken{Token: token, CreatedAt: time.Now()}
	if err := r.db.Create(&blacklist).Error; err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

// IsTokenBlacklisted memeriksa apakah token ada dalam tabel blacklist
func (r *UserRepository) IsTokenBlacklisted(token string) bool {
	var blacklist model.BlacklistedToken
	err := r.db.Where("token = ?", token).First(&blacklist).Error
	return err == nil // Jika ditemukan, token di-blacklist
}
