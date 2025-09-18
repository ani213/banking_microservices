package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

var validationMessages = map[string]string{
	"Email.required":    "Email is required",
	"Email.email":       "Please enter a valid email address",
	"Password.required": "Password is required",
	"Password.min":      "Password must be at least 8 characters",
	"FullName.required": "Full name is required",
}

func Error(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func Success(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func ValidationError(w http.ResponseWriter, err error, code int) {
	var messages []map[string]string
	for _, e := range err.(validator.ValidationErrors) {
		key := e.Field() + "." + e.Tag()
		msg, ok := validationMessages[key]
		if !ok {
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
