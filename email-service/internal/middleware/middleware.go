package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/ani213/email-service/internal/config"
	"github.com/ani213/email-service/utils"
)

type User struct {
	UserID   int64
	Email    string
	Roles    []int64
	FullName string
}

type validateResponse struct {
	Valid bool   `json:"vaild"`
	User  User   `json:"user"`
	Error string `json:"error"`
}

type ContextKey string

const UserContextKey ContextKey = "user"

func Authenticate(config *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("line 25", authHeader)
				utils.UnAutherizedError(w)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Println("lin3 32", parts[0])
				utils.UnAutherizedError(w)
				return
			}
			token := parts[1]

			autRequest, err := http.NewRequest("GET", config.AuthService+"/auth/validate-token", nil)
			if err != nil {
				http.Error(w, "Auth server unavailabel", http.StatusInternalServerError)
				return
			}
			autRequest.Header.Set("Content-Type", "application/json")
			autRequest.Header.Set("Authorization", "Bearer "+token)
			client := &http.Client{}
			resp, err := client.Do(autRequest)
			if err != nil {
				log.Println(err.Error(), "Error main things")
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				utils.UnAutherizedError(w)
				return
			}
			var result validateResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || !result.Valid {
				utils.UnAutherizedError(w)
				return
			}
			ctx := context.WithValue(r.Context(), UserContextKey, result)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
