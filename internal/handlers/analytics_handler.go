package handlers

import (
	"net/http"
	"strconv"

	"Zira/internal/services"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// AnalyticsHandler обрабатывает запросы, связанные с аналитикой
type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
	logger           *logrus.Logger
}

// NewAnalyticsHandler создает новый AnalyticsHandler
func NewAnalyticsHandler(analyticsService *services.AnalyticsService, logger *logrus.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService, logger: logger}
}

// GetUserAnalytics возвращает статистику по доходам/расходам (заглушка)
func (h *AnalyticsHandler) GetUserAnalytics(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]interface{}{"income": 0, "expense": 0}, http.StatusOK)
}

// PredictBalance прогнозирует баланс на N дней (заглушка)
func (h *AnalyticsHandler) PredictBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, _ := strconv.ParseUint(vars["accountId"], 10, 64)
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days <= 0 {
		days = 30
	}
	respondJSON(w, map[string]interface{}{"account_id": accountID, "predicted_balance": 0, "days": days}, http.StatusOK)
}
