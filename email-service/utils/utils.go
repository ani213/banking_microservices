package utils

import (
	"encoding/json"
	"net/http"
)

func UnAutherizedError(w http.ResponseWriter) {

	w.WriteHeader(http.StatusUnauthorized)
	msg := "Unautherized access"

	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func CustomError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
