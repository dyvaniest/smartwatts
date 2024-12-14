package repository

import (
	"a21hc3NpZ25tZW50/model"
	"fmt"
	"sync"
	"time"
)

type UserRepository struct {
	mu                sync.Mutex
	users             []model.User
	blacklistedTokens map[string]bool
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:             []model.User{},
		blacklistedTokens: make(map[string]bool),
	}
}

func (r *UserRepository) AddUser(user model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.users {
		if u.Email == user.Email {
			return fmt.Errorf("email already registered")
		}
	}

	user.ID = len(r.users) + 1
	user.CreatedAt = time.Now()
	r.users = append(r.users, user)
	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, user := range r.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

// BlacklistToken menambahkan token ke daftar blacklist
func (r *UserRepository) BlacklistToken(token string) {
	r.blacklistedTokens[token] = true
}

// IsTokenBlacklisted memeriksa apakah token ada dalam blacklist
func (r *UserRepository) IsTokenBlacklisted(token string) bool {
	return r.blacklistedTokens[token]
}
