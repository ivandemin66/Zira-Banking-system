package models

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID              uint    `gorm:"primaryKey;autoIncrement"`
	FromAccountID   uint    `gorm:"not null;index"`
	ToAccountID     uint    `gorm:"not null;index"`
	Amount          float64 `gorm:"type:numeric;not null"`
	Currency        string  `gorm:"size:3;not null"`
	TransactionType string  `gorm:"size:50;not null"` // e.g. "transfer", "deposit"
	CreatedAt       time.Time
}

// TransferRequest описывает структуру запроса на перевод средств между счетами
// Используется для валидации и сериализации входных данных
// swagger:model
type TransferRequest struct {
	FromAccountID uint    `json:"from_account_id" validate:"required"`
	ToAccountID   uint    `json:"to_account_id" validate:"required,nefield=FromAccountID"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Currency      string  `json:"currency" validate:"required,eq=RUB"`
}

// Validate выполняет базовую валидацию структуры TransferRequest
func (t *TransferRequest) Validate() error {
	if t.FromAccountID == t.ToAccountID {
		return fmt.Errorf("Счета отправителя и получателя не должны совпадать")
	}
	if t.Amount <= 0 {
		return fmt.Errorf("Сумма перевода должна быть больше нуля")
	}
	if t.Currency != "RUB" {
		return fmt.Errorf("Поддерживается только валюта RUB")
	}
	return nil
}
