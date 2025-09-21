package auth

import (
	"github.com/gorilla/mux"
)

func PublicRoutes(r *mux.Router, h *Handler) {
	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")

}

func PrivateRoutes(r *mux.Router, h *Handler) {
	r.HandleFunc("/users", h.GetUsers).Methods("GET")
	r.HandleFunc("/context", h.GetContext).Methods("GET")
	r.HandleFunc("/validate-token", h.ValidateToken).Methods("GET")

}
