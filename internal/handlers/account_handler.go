package handlers

import (
	"Zira/internal/models"
	"Zira/internal/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// AccountService определяет интерфейс для работы со счетами
// Интерфейс должен быть реализован в пакете services
// type AccountService interface { ... }

type AccountHandler struct {
	accountService services.AccountService // Используем интерфейс из пакета services
	logger         *logrus.Logger
}

func NewAccountHandler(accountService services.AccountService, logger *logrus.Logger) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		logger:         logger,
	}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	var req models.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Ошибка декодирования JSON")
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	account, err := h.accountService.CreateAccount(userID, req.Type)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка создания счета")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.respondJSON(w, account, http.StatusCreated)
}

func (h *AccountHandler) GetUserAccounts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	accounts, err := h.accountService.GetUserAccounts(userID)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка получения счетов пользователя")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, accounts, http.StatusOK)
}

func (h *AccountHandler) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	vars := mux.Vars(r)
	accountID, err := strconv.ParseInt(vars["accountId"], 10, 64)
	if err != nil {
		h.logger.WithError(err).Error("Некорректный ID счета")
		http.Error(w, "Некорректный ID счета", http.StatusBadRequest)
		return
	}

	account, err := h.accountService.GetAccountByID(accountID)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка получения счета")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Приведение account к *models.Account
	acc, ok := account.(*models.Account)
	if !ok {
		h.logger.Error("Ошибка типа account")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Проверка прав доступа
	if acc.UserID != uint(userID) {
		h.logger.WithFields(logrus.Fields{
			"accountID": accountID,
			"userID":    userID,
		}).Warn("Попытка доступа к чужому счету")
		http.Error(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	h.respondJSON(w, acc, http.StatusOK)
}

func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	var req models.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Ошибка декодирования JSON")
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		h.logger.WithError(err).Error("Ошибка валидации запроса перевода")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка владения счетом отправителя
	fromAccount, err := h.accountService.GetAccountByID(int64(req.FromAccountID))
	if err != nil {
		h.logger.WithError(err).Error("Счет отправителя не найден")
		http.Error(w, "Счет отправителя не найден", http.StatusNotFound)
		return
	}

	fromAcc, ok := fromAccount.(*models.Account)
	if !ok {
		h.logger.Error("Ошибка типа fromAccount")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	if fromAcc.UserID != uint(userID) {
		h.logger.WithFields(logrus.Fields{
			"accountID": req.FromAccountID,
			"userID":    userID,
		}).Warn("Попытка перевода с чужого счета")
		http.Error(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	transaction, err := h.accountService.Transfer(&req)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка перевода средств")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.respondJSON(w, transaction, http.StatusOK)
}

func (h *AccountHandler) respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.WithError(err).Error("Ошибка кодирования JSON")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}
