package main

import (
	"Zira/internal/config"
	"Zira/internal/handlers"
	"Zira/internal/middleware"
	"Zira/internal/repository"
	"Zira/internal/routes"
	"Zira/internal/services"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// Загрузка конфигурации
	cfg := config.LoadConfig("config/config.yaml")

	// Настройка логгера
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	// Подключение к БД через GORM
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Pass, cfg.Database.Name,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.WithError(err).Fatal("Не удалось подключиться к базе данных")
	}

	// Настройка репозиториев
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	cardRepo := repository.NewCardRepository(db)
	creditRepo := repository.NewCreditRepository(db)
	paymentScheduleRepo := repository.NewPaymentScheduleRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Настройка сервисов
	userService := services.NewUserService(userRepo, cfg.JWT.Host) // или cfg.JWT.Secret, если есть
	accountService := services.NewAccountService(accountRepo)
	cardService := services.NewCardService(cardRepo, "PGP_PUBLIC_KEY", []byte("hmac_secret"))
	emailService := services.NewEmailService("smtp.example.com", 587, "noreply@example.com", "password")
	creditService := services.NewCreditService(creditRepo, paymentScheduleRepo, accountRepo, emailService)
	analyticsService := services.NewAnalyticsService(transactionRepo, accountRepo, creditRepo, paymentScheduleRepo)

	// Настройка middleware
	jwtAuth := middleware.Auth(cfg.JWT.Host) // или cfg.JWT.Secret

	// Настройка обработчиков
	userHandler := handlers.NewUserHandler(userService, logger)
	accountHandler := handlers.NewAccountHandler(accountService, logger)
	cardHandler := handlers.NewCardHandler(cardService, logger)
	creditHandler := handlers.NewCreditHandler(creditService, logger)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, logger)

	// Настройка маршрутизатора
	router := mux.NewRouter()
	// Публичные маршруты
	router.HandleFunc("/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")
	// Защищённые маршруты
	router.Handle("/accounts", jwtAuth(http.HandlerFunc(accountHandler.CreateAccount))).Methods("POST")
	// ... добавьте остальные защищённые маршруты по аналогии

	// Приветственный обработчик для корня
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Добро пожаловать в Zira Banking System API! Используйте /register, /login и другие эндпоинты. Подробнее см. README.md."))
	}).Methods("GET")

	routes.SetupRoutes(router, userHandler, accountHandler, cardHandler, creditHandler, analyticsHandler)

	logger.Info("Сервер запущен на порту 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.WithError(err).Fatal("Ошибка запуска сервера")
	}
}
