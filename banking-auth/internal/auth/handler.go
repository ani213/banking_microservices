package auth

import (
	"encoding/json"
	"net/http"

	"github.com/ani213/banking-auth/util"
	"github.com/go-playground/validator"
)

type Handler struct {
	svc *Service
}

var validate = validator.New()

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	json.NewDecoder(r.Body).Decode(&req)
	err := validate.Struct(req)
	if err != nil {
		util.ValidationError(w, err, http.StatusBadRequest)
		return
	}
	if err := h.svc.Register(&req); err != nil {
		util.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.Success(w, "user registered successfully", http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		util.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}
