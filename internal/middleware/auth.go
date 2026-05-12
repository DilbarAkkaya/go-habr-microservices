package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/DilbarAkkaya/go-habr-microservices/internal/utils"
)

type contextKey string

const userContextKey contextKey = "user"

type AuthUser struct {
	ID    int
	Email string
}

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Middleware called for %s %s", r.Method, r.URL.Path)
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "missing auth header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseJWT(parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		user := AuthUser{
			ID:    claims.UserID,
			Email: claims.Email,
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) (AuthUser, bool) {
	user, ok := ctx.Value(userContextKey).(AuthUser)
	return user, ok
}
