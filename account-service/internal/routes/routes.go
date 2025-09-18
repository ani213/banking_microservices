package routes

import (
	"github.com/ani213/account-service/internal/account"
	"github.com/gorilla/mux"
)

func Routes(r *mux.Router, h *account.Handler) {
	r.HandleFunc("/accounts", h.CreateAccount).Methods("POST")
	r.HandleFunc("/accounts/{id}", h.GetAccount).Methods("GET")
	r.HandleFunc("/accounts/{id}/deposit", h.Deposit).Methods("POST")
	r.HandleFunc("/accounts/{id}/withdraw", h.Withdraw).Methods("POST")
}
