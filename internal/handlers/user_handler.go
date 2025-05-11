package handlers

import (
	"Zira/internal/models"
	"Zira/internal/services"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// UserService определяет интерфейс для работы с пользователями
// Интерфейс должен быть реализован в пакете services
// type UserService interface { ... }

type UserHandler struct {
	userService services.UserService // Используем интерфейс из пакета services
	logger      *logrus.Logger
}

func NewUserHandler(userService services.UserService, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Ошибка декодирования JSON")
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		h.logger.WithError(err).Error("Ошибка валидации данных")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	auth, err := h.userService.Register(&req)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка регистрации пользователя")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.respondJSON(w, auth, http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Ошибка декодирования JSON")
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	auth, err := h.userService.Login(&req)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка аутентификации пользователя")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.respondJSON(w, auth, http.StatusOK)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка получения профиля пользователя")
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	h.respondJSON(w, user, http.StatusOK)
}

func (h *UserHandler) respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.WithError(err).Error("Ошибка кодирования JSON")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}
