package account

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ani213/account-service/internal/util"
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

func GetContexValue(r *http.Request) ContextValue {
	val := r.Context().Value(UserContextKey)
	user, _ := val.(ContextValue)
	return user
}

var validate = validator.New()

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var acc Account
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(acc); err != nil {
		util.ValidationError(w, err, http.StatusBadRequest)
		// utils.CustomError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.CreateAccount(r.Context(), &acc); err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		util.CustomError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	contextUser := GetContexValue(r)
	go h.service.SendEmail(contextUser.Email, "Account Creation", "Your Account number:-"+acc.AccountNumber+" is successfully created", r)
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

func (h *Handler) GetAccountsByUserID(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	accounts, err := h.service.GetAccountsByUserID(userID)
	if err != nil {
		util.CustomError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"accounts": accounts})
}

func (h *Handler) GetAllUserWithAccounts(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUserWithAccounts()
	if err != nil {
		util.CustomError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"users": users})
}
