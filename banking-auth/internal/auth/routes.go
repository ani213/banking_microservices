package auth

import (
	"github.com/gorilla/mux"
)

func Routes(r *mux.Router, h *Handler) {
	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/validate-token", h.ValidateToken).Methods("POST")
	r.HandleFunc("/users", h.GetUsers).Methods("GET")
}
