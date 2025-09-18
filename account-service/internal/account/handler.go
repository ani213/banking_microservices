package account

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ani213/account-service/internal/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

var validate = validator.New()

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var acc Account
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(acc); err != nil {
		utils.ValidationError(w, err, http.StatusBadRequest)
		// utils.CustomError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.CreateAccount(r.Context(), &acc); err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		utils.CustomError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	contextValue := utils.GetContextValue(r)
	log.Println(contextValue, "context value")
	// h.service.SendEmail()
	json.NewEncoder(w).Encode(acc)
}

func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	acc, err := h.service.GetAccount(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(acc)
}

func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var payload struct {
		Amount string `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	amount, _ := decimal.NewFromString(payload.Amount)
	if err := h.service.Deposit(r.Context(), id, amount); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var payload struct {
		Amount string `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	amount, _ := decimal.NewFromString(payload.Amount)
	if err := h.service.Withdraw(r.Context(), id, amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
