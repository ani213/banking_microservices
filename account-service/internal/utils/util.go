package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ani213/account-service/internal/middleware"
	"github.com/go-playground/validator"
)

func CustomError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func ValidationError(w http.ResponseWriter, err error, code int) {
	var messages []map[string]string
	for _, e := range err.(validator.ValidationErrors) {
		var msg string
		switch e.Tag() {
		case "required":
			msg = fmt.Sprintf("%s is required", e.Field())
		default:
			msg = fmt.Sprintf("%s is invalid", e.Field())
		}

		messages = append(messages, map[string]string{
			"field":   e.Field(),
			"message": msg,
		})

	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{"errors": messages})
}

func GetToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	log.Println(authHeader, "Auth header>>>")
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 {
		return parts[1]
	} else {
		return ""
	}
}

func GetContextValue(r *http.Request) string {
	user := r.Context().Value(middleware.UserIDKey)
	log.Println(user, "user")
	return "ani"
}
