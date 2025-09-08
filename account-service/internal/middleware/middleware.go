package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ani213/auth-service/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// var jwtSecret = []byte("super-secret-key") // should come from ENV

type contextKey string

const UserIDKey contextKey = "userID"

// JWTMiddleware validates JWT from Auth Service
// func JWTMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			AuthError(w)
// 			return
// 		}

// 		parts := strings.Split(authHeader, " ")
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			AuthError(w)

// 			// http.Error(w, "invalid token format", http.StatusUnauthorized)
// 			return
// 		}

// 		tokenStr := parts[1]
// 		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
// 			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, jwt.ErrInvalidKey
// 			}
// 			fmt.Println(jwtSecret, "secret from env")
// 			return jwtSecret, nil
// 		})

// 		if err != nil || !token.Valid {
// 			AuthError(w)

// 			// http.Error(w, "invalid token", http.StatusUnauthorized)
// 			return
// 		}

// 		claims := token.Claims.(jwt.MapClaims)
// 		userID, err := getUserIDFromClaims(claims)
// 		if err != nil {
// 			AuthError(w)

// 			// http.Error(w, "invalid user_id in token", http.StatusUnauthorized)
// 			return
// 		}

// 		// attach userID to request context
// 		ctx := context.WithValue(r.Context(), UserIDKey, userID)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func JWTMiddleware(cnf *config.Config) func(http.Handler) http.Handler {
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
				return []byte(cnf.JwtSecret), nil
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
