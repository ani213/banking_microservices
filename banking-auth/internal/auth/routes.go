package auth

import (
	"github.com/gorilla/mux"
)

func Routes(r *mux.Router, h *Handler) {
	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
}
