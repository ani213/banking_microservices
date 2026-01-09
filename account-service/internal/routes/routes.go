package routes

import (
	"github.com/ani213/account-service/internal/account"
	"github.com/gorilla/mux"
)

func Routes(r *mux.Router, h *account.Handler) {

	r.HandleFunc("", h.CreateAccount).Methods("POST")
	r.HandleFunc("/{id}", h.GetAccount).Methods("GET")
	r.HandleFunc("/{id}/deposit", h.Deposit).Methods("POST")
	r.HandleFunc("/{id}/withdraw", h.Withdraw).Methods("POST")
	r.HandleFunc("", h.GetAllUserWithAccounts).Methods("GET")
	r.HandleFunc("/user/{user_id}", h.GetAccountsByUserID).Methods("GET")
}
