package service

import (
	"a21hc3NpZ25tZW50/model"
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type UserService struct {
	userRepo *repository.UserRepository
}

// GenerateToken menghasilkan JWT untuk pengguna
func GenerateToken(email string) (string, error) {
	expirationTime := time.Now().Add(20 * time.Minute)
	claims := &model.Claim{
		Email: email,
		Claims: jwt.MapClaims{
			"exp": expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) RegisterUser(user model.User) error {
	return s.userRepo.AddUser(user)
}

func (s *UserService) LoginUser(req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if user.Password != req.Password {
		return nil, fmt.Errorf("invalid password")
	}

	token, _ := GenerateToken(user.Email)
	return &model.LoginResponse{
		Message: "Login successfullty",
		Token:   token,
	}, nil
}

// LogoutUser mem-blacklist token yang di-logout
func (s *UserService) LogoutUser(token string) {
	s.userRepo.BlacklistToken(token)
}

// IsLoggedOut memeriksa apakah token di-blacklist
func (s *UserService) IsLoggedOut(token string) bool {
	return s.userRepo.IsTokenBlacklisted(token)
}
