package auth

import (
	"github.com/ani213/banking-auth/internal/middleware"
	"github.com/gorilla/mux"
)

func PublicRoutes(r *mux.Router, h *Handler) {
	api := r.PathPrefix("/").Subrouter()
	api.Use(middleware.JSONMiddleware)
	api.HandleFunc("/register", h.Register).Methods("POST")
	api.HandleFunc("/login", h.Login).Methods("POST")

}

func PrivateRoutes(r *mux.Router, h *Handler) {
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JSONMiddleware)
	api.HandleFunc("/context", h.GetContext).Methods("GET")
	api.HandleFunc("/users", h.GetUsers).Methods("GET")
	api.HandleFunc("/validate-token", h.ValidateToken).Methods("GET")
	api.HandleFunc("/add-role", h.AddUserRole).Methods("PUT")
}
