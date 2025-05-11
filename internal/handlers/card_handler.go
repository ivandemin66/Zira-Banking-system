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

// CardHandler обрабатывает запросы, связанные с банковскими картами
type CardHandler struct {
	cardService *services.CardService
	logger      *logrus.Logger
}

// NewCardHandler создает новый CardHandler
func NewCardHandler(cardService *services.CardService, logger *logrus.Logger) *CardHandler {
	return &CardHandler{cardService: cardService, logger: logger}
}

// CreateCard выпускает новую карту
func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID   uint   `json:"account_id"`
		ExpiryMonth int    `json:"expiry_month"`
		ExpiryYear  int    `json:"expiry_year"`
		CVV         string `json:"cvv"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Ошибка декодирования JSON")
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}
	card, err := h.cardService.GenerateCard(req.AccountID, req.ExpiryMonth, req.ExpiryYear, req.CVV)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка выпуска карты")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, card, http.StatusCreated)
}

// GetUserCards возвращает список карт пользователя
func (h *CardHandler) GetUserCards(w http.ResponseWriter, r *http.Request) {
	// Здесь должна быть реализация поиска карт по userID (заглушка)
	respondJSON(w, []models.Card{}, http.StatusOK)
}

// GetCardByID возвращает карту по ID
func (h *CardHandler) GetCardByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID, err := strconv.ParseUint(vars["cardId"], 10, 64)
	if err != nil {
		http.Error(w, "Некорректный ID карты", http.StatusBadRequest)
		return
	}
	// Здесь должна быть реализация поиска и расшифровки карты (заглушка)
	respondJSON(w, &models.Card{ID: uint(cardID)}, http.StatusOK)
}

// respondJSON формирует JSON-ответ
func respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
