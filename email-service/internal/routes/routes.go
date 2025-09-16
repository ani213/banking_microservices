package routes

import (
	"github.com/ani213/email-service/internal/email"
	"github.com/gorilla/mux"
)

func Routes(r *mux.Router, h *email.Handler) {
	r.HandleFunc("/send-email", h.SendEmail).Methods("POST")
}
