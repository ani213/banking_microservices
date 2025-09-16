package middleware

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/ani213/email-service/internal/config"
	"github.com/ani213/email-service/utils"
)

type validateResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"userId,omitempty"`
	Error  string `json:"error,omitempty"`
}

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

			reqBody, _ := json.Marshal(map[string]string{"token": token})
			resp, err := http.Post(config.AuthService+"/validate-token", "application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				http.Error(w, "Auth server unavailabel", http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Println("not ok status", resp.StatusCode)
				utils.UnAutherizedError(w)
				return
			}
			var result validateResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || !result.Valid {
				utils.UnAutherizedError(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
