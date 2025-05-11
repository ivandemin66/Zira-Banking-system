package models

import (
	"errors"
	"regexp"
	"time"
	"unicode"
)

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"size:100;uniqueIndex;not null"`
	Password  string `gorm:"size:255;not null"` // захэшированный bcrypt
	CreatedAt time.Time
	UpdatedAt time.Time

	Accounts []Account `gorm:"foreignKey:UserID"`
}

// RegisterRequest описывает структуру запроса на регистрацию пользователя
// Используется для валидации и сериализации входных данных
// Валидация email и пароля должна быть реализована в сервисе
// swagger:model
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

// LoginRequest описывает структуру запроса на аутентификацию пользователя
// Используется для валидации и сериализации входных данных
// swagger:model
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

// Validate выполняет полную валидацию полей регистрации
func (r *RegisterRequest) Validate() error {
	if len(r.Name) < 2 || len(r.Name) > 100 {
		return errors.New("Имя должно быть от 2 до 100 символов")
	}
	// Email: базовая проверка через regexp
	re := regexp.MustCompile(`^[\w._%+-]+@[\w.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(r.Email) {
		return errors.New("Некорректный email")
	}
	if len(r.Password) < 8 || len(r.Password) > 64 {
		return errors.New("Пароль должен быть от 8 до 64 символов")
	}
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, c := range r.Password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return errors.New("Пароль должен содержать строчные, прописные буквы, цифры и спецсимволы")
	}
	return nil
}
