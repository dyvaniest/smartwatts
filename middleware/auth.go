package middleware

import (
	"net/http"
	"os"

	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"a21hc3NpZ25tZW50/service"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

var userRepo = repository.NewUserRepository()
var userService = service.NewUserService(userRepo)

var sessionService = service.NewSessionService("my-secret-key", "auth-session", "auth-token")

// AuthenticationMiddleware memvalidasi JWT token pada permintaan
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := sessionService.GetTokenFromSession(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Periksa apakah token di-blacklist
		if userService.IsLoggedOut(tokenStr) {
			http.Error(w, "Token is invalid or logged out", http.StatusUnauthorized)
			return
		}

		// Validasi token
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or Expired Token", http.StatusUnauthorized)
			return
		}

		// Jika valid, tambahkan informasi pengguna ke context
		email := (*claims)["email"].(string)
		r.Header.Set("Authenticated-Email", email)

		// Lanjutkan ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}
