package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Zira/internal/models"
	"Zira/internal/services"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// CreditHandler обрабатывает запросы, связанные с кредитами
type CreditHandler struct {
	creditService *services.CreditService
	logger        *logrus.Logger
}

// NewCreditHandler создает новый CreditHandler
func NewCreditHandler(creditService *services.CreditService, logger *logrus.Logger) *CreditHandler {
	return &CreditHandler{creditService: creditService, logger: logger}
}

// ApplyForCredit оформляет новый кредит
func (h *CreditHandler) ApplyForCredit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID uint    `json:"account_id"`
		Amount    float64 `json:"amount"`
		Months    int     `json:"months"`
		Rate      float64 `json:"rate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Ошибка декодирования JSON")
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}
	credit, err := h.creditService.ApplyForCredit(req.AccountID, req.Amount, req.Months, req.Rate)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка оформления кредита")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, credit, http.StatusCreated)
}

// GetUserCredits возвращает список кредитов пользователя (заглушка)
func (h *CreditHandler) GetUserCredits(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, []models.Credit{}, http.StatusOK)
}

// GetCreditByID возвращает детали кредита (заглушка)
func (h *CreditHandler) GetCreditByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	creditID, err := strconv.ParseUint(vars["creditId"], 10, 64)
	if err != nil {
		http.Error(w, "Некорректный ID кредита", http.StatusBadRequest)
		return
	}
	respondJSON(w, &models.Credit{ID: uint(creditID)}, http.StatusOK)
}

// GetPaymentSchedule возвращает график платежей по кредиту (заглушка)
func (h *CreditHandler) GetPaymentSchedule(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, []models.PaymentSchedule{}, http.StatusOK)
}
