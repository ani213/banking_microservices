package jwtutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-secret-key") // in prod, use env vars

type contextKey string

const UserContextKey contextKey = "user"

type ContextValue struct {
	UserId   int64
	Email    string
	Roles    []int64
	FullName string
}

func GenerateToken(userID int64, email string, fullName string, roles []int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"email":    email,
		"fullName": fullName,
		"roles":    roles,
	}
	fmt.Println(roles, "roles")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenStr string) (ContextValue, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return ContextValue{}, err
	}

	claims := token.Claims.(jwt.MapClaims)
	user, err := getUserFromClaim(claims)
	if err != nil {
		return ContextValue{}, err
	}
	return user, nil
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
			user, err := getUserFromClaim(claims)
			if err != nil {
				AuthError(w)
				return
			}
			fmt.Println(user, "user values")
			// attach user to request context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extract userID safely from claims
func getUserFromClaim(claims jwt.MapClaims) (ContextValue, error) {
	var user ContextValue
	if val, ok := claims["user_id"].(float64); ok {
		user.UserId = int64(val)
		// return int64(val), nil
	}
	if val, ok := claims["user_id"].(string); ok {
		// try to parse string to int64
		parsed, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return user, fmt.Errorf("invalid user_id format")
		}
		user.UserId = parsed
	}
	if email, ok := claims["email"].(string); ok {
		user.Email = string(email)
	}
	if rawRoles, ok := claims["roles"].([]interface{}); ok {
		var roles []int64
		for _, r := range rawRoles {
			if roleFloat, ok := r.(float64); ok {
				roles = append(roles, int64(roleFloat))
			}
		}
		user.Roles = roles
	}
	if fullName, ok := claims["fullName"].(string); ok {
		user.FullName = fullName
	}

	return user, nil
}

func AuthError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized access"})
}

func GetContextValue(r *http.Request) (ContextValue, error) {
	contextValue, ok := r.Context().Value(UserContextKey).(ContextValue)
	if !ok {
		return ContextValue{}, errors.New("get context value error")
	}
	user := ContextValue{
		UserId:   contextValue.UserId,
		Email:    contextValue.Email,
		Roles:    contextValue.Roles,
		FullName: contextValue.FullName,
	}
	return user, nil
}
