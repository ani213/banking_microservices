package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ani213/banking-auth/pkg/jwtutil"
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

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := h.svc.GetUsers()
	if err != nil {
		util.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string][]ResponsGetUser{
		"users": users,
	})
}

func (h *Handler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	user, err := h.svc.ValidateToken(r)
	if err != nil {
		util.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"user": user, "valid": true})
}

func (h *Handler) GetContext(w http.ResponseWriter, r *http.Request) {
	user, err := jwtutil.GetContextValue(r)
	if err != nil {
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) AddUserRole(w http.ResponseWriter, r *http.Request) {
	var requestBody UserRoleRequestBody
	json.NewDecoder(r.Body).Decode(&requestBody)
	err := validate.Struct(requestBody)
	if err != nil {
		util.ValidationError(w, err, http.StatusBadRequest)
		return
	}
	fmt.Println(requestBody, "request body")
	err = h.svc.AddRoles(requestBody.UserID, requestBody.Roles)
	if err != nil {
		util.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"msg": "Successfully role added"})
}
