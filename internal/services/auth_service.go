package services

import (
	"fmt"

	"Zira/internal/models"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// AuthService реализует бизнес-логику аутентификации пользователей
// Здесь будут методы регистрации, логина и т.д.
type AuthService struct {
	DB *gorm.DB
}

// NewAuthService создает новый сервис аутентификации
func NewAuthService(dsn string) (*AuthService, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}
	return &AuthService{DB: db}, nil
}

// UserService определяет интерфейс для работы с пользователями
// Используется в хендлерах для инверсии зависимостей
// Реализация UserService находится в user_service.go
type UserService interface {
	Register(req *models.RegisterRequest) (interface{}, error)
	Login(req *models.LoginRequest) (interface{}, error)
	GetUserByID(id int64) (interface{}, error)
}

// AccountService определяет интерфейс для работы со счетами
// Используется в хендлерах для инверсии зависимостей
type AccountService interface {
	CreateAccount(userID int64, accountType string) (interface{}, error)
	GetUserAccounts(userID int64) (interface{}, error)
	GetAccountByID(accountID int64) (interface{}, error)
	Transfer(req *models.TransferRequest) (interface{}, error)
}
