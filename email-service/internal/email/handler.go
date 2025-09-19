package email

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var reqBody EmailRequest
	json.NewDecoder(r.Body).Decode(&reqBody)
	msg, _ := h.service.SendEmail(&reqBody)
	fmt.Println(msg)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"msg": msg})
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"msg": "test OK"})
}
