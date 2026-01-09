package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ani213/account-service/internal/account"
	"github.com/ani213/account-service/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// var jwtSecret = []byte("super-secret-key") // should come from ENV

// type contextKey string

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

			claims, ok := token.Claims.(jwt.MapClaims)

			user, err := getUserIDFromClaims(claims, ok)

			if err != nil {
				AuthError(w)
				return
			}
			// attach userID to request context
			ctx := context.WithValue(r.Context(), account.UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extract userID safely from claims
func getUserIDFromClaims(claims jwt.MapClaims, ok bool) (account.ContextValue, error) {
	var contextValue account.ContextValue
	if !ok {
		return account.ContextValue{}, fmt.Errorf("not getting context value")
	}
	if val, ok := claims["user_id"].(float64); ok {
		contextValue.User_id = int64(val)
	}
	if val, ok := claims["user_id"].(string); ok {
		// try to parse string to int64
		parsed, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return contextValue, fmt.Errorf("invalid user_id format")
		}
		contextValue.User_id = parsed
	}
	// contextValue.User_id = (claims["user_id"].(int64))
	contextValue.Email = (claims["email"]).(string)
	contextValue.FullName = claims["fullName"].(string)

	return contextValue, nil
}

func AuthError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized access"})
}

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "applicatio/json")
		next.ServeHTTP(w, r)
	})
}
