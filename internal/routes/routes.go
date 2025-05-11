package routes

import (
	"Zira/internal/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(
	r *mux.Router,
	userHandler *handlers.UserHandler,
	accountHandler *handlers.AccountHandler,
	cardHandler *handlers.CardHandler,
	creditHandler *handlers.CreditHandler,
	analyticsHandler *handlers.AnalyticsHandler,
	// authMiddleware *middleware.AuthMiddleware,
) {
	// Публичные маршруты
	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	// Защищенные маршруты
	protected := r.PathPrefix("/").Subrouter()
	// protected.Use(authMiddleware.Authenticate)

	// Профиль пользователя
	protected.HandleFunc("/profile", userHandler.GetProfile).Methods("GET")

	// Счета
	protected.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	protected.HandleFunc("/accounts", accountHandler.GetUserAccounts).Methods("GET")
	protected.HandleFunc("/accounts/{accountId}", accountHandler.GetAccountByID).Methods("GET")
	protected.HandleFunc("/transfer", accountHandler.Transfer).Methods("POST")

	// Карты
	protected.HandleFunc("/cards", cardHandler.CreateCard).Methods("POST")
	protected.HandleFunc("/cards", cardHandler.GetUserCards).Methods("GET")
	protected.HandleFunc("/cards/{cardId}", cardHandler.GetCardByID).Methods("GET")

	// Транзакции
	// protected.HandleFunc("/transactions", transactionHandler.GetUserTransactions).Methods("GET")
	// protected.HandleFunc("/accounts/{accountId}/transactions", transactionHandler.GetAccountTransactions).Methods("GET")

	// Кредиты
	protected.HandleFunc("/credits", creditHandler.ApplyForCredit).Methods("POST")
	protected.HandleFunc("/credits", creditHandler.GetUserCredits).Methods("GET")
	protected.HandleFunc("/credits/{creditId}", creditHandler.GetCreditByID).Methods("GET")
	protected.HandleFunc("/credits/{creditId}/schedule", creditHandler.GetPaymentSchedule).Methods("GET")

	// Аналитика
	protected.HandleFunc("/analytics", analyticsHandler.GetUserAnalytics).Methods("GET")
	protected.HandleFunc("/accounts/{accountId}/predict", analyticsHandler.PredictBalance).Methods("GET")
}
