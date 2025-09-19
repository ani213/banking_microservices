package jwtutil

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-secret-key") // in prod, use env vars

type contextKey string

const UserIDKey contextKey = "userID"

func GenerateToken(userID string, email string, fullName string, roles []int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"email":    email,
		"fullName": fullName,
		"roles":    roles,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenStr string) (int64, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)
	userID, err := getUserIDFromClaims(claims)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func JWTMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				AuthError(w)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				AuthError(w)
				return
			}

			tokenStr := parts[1]
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrInvalidKey
				}
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				AuthError(w)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			userID, err := getUserIDFromClaims(claims)
			if err != nil {
				AuthError(w)
				return
			}

			// attach userID to request context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extract userID safely from claims
func getUserIDFromClaims(claims jwt.MapClaims) (int64, error) {
	if val, ok := claims["user_id"].(float64); ok {
		return int64(val), nil
	}
	if val, ok := claims["user_id"].(string); ok {
		// try to parse string to int64
		parsed, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid user_id format")
		}
		return parsed, nil
	}
	return 0, fmt.Errorf("user_id not found or invalid type")
}

func AuthError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized access"})
}
