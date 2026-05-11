package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/finalops/auth-service/pkg/jwtutil"
	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtutil.ValidateToken(tokenStr, secret)
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) uuid.UUID {
	v := r.Context().Value(UserIDKey)
	if id, ok := v.(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}
