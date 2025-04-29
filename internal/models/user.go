package models

import (
	"errors"
	"regexp"
	"time"
)

type user struct {
	ID           int64     `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  user   `json:"user"`
}

// Validate проверяет корректность данных регистрации
func (r *RegisterRequest) Validate() error {
	if len(r.Username) < 3 {
		return errors.New("имя пользователя не должно содержать менее 3 символов")
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(r.Email) {
		return errors.New("некорректный формат email")
	}
	if len(r.Password) < 8 {
		return errors.New("пароль должен содержать не менее 8 символов")
	}
	return nil
}
