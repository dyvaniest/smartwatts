package middleware

import (
	"net/http"
	"os"

	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"a21hc3NpZ25tZW50/service"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

var sessionRepo = repository.NewSessionRepo(&gorm.DB{})
var sessionService = service.NewSessionService(sessionRepo)

// AuthenticationMiddleware memvalidasi JWT token pada permintaan
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := sessionService.GetTokenFromSession(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		session, err := sessionService.ValidateSession(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
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

		r.Header.Set("Authenticated-Email", session.Email)

		// Lanjutkan ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}
