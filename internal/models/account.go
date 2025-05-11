package models

import "time"

type Account struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	UserID    uint    `gorm:"not null;index"`
	Balance   float64 `gorm:"type:numeric;not null;default:0"`
	Currency  string  `gorm:"size:3;not null"` // ISO-код валюты, например "RUB"
	CreatedAt time.Time
	UpdatedAt time.Time

	Cards []Card `gorm:"foreignKey:AccountID"`
}

// CreateAccountRequest описывает структуру запроса на создание счета
// Используется для валидации и сериализации входных данных
// swagger:model
type CreateAccountRequest struct {
	Type     string `json:"type" validate:"required,oneof=debit credit savings"`
	Currency string `json:"currency" validate:"required,eq=RUB"`
}
