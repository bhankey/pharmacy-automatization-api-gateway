package middleware

import (
	"context"
	"net/http"

	"github.com/bhankey/go-utils/pkg/logger"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
	"github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
	jwtKey string

	logger logger.Logger
}

func NewAuthMiddleware(log logger.Logger, jwtKey string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtKey: jwtKey,
		logger: log,
	}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte("failed to authorize user")); err != nil {
				m.logger.Warn("failed to write to user", err)
			}

			return
		}

		claim := entities.Claims{}
		token, err := jwt.ParseWithClaims(authHeader, &claim, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.jwtKey), nil
		})
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte("failed to authorize user")); err != nil {
				m.logger.Warn("failed to write to user", err)
			}

			return
		}

		// Better to move this to microservice and get only
		ctx := r.Context()
		ctx = context.WithValue(ctx, entities.UserID, claim.UserID)
		ctx = context.WithValue(ctx, entities.Email, claim.Email)
		ctx = context.WithValue(ctx, entities.PharmacyID, claim.PharmacyID)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
