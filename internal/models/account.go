package models

import (
	"errors"
	"time"
)

type Account struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Number    string    `json:"number" db:"number"`
	Balance   float64   `json:"balance" db:"balance"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateAccountRequest struct {
	Type string `json:"type"`
}

type TransferRequest struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
}

func (r *TransferRequest) Validate() error {
	if r.FromAccountID <= 0 {
		return errors.New("некорректный идентификатор счета отправителя")
	}
	if r.ToAccountID <= 0 {
		return errors.New("некорректный идентификатор счета получателя")
	}
	if r.Amount <= 0 {
		return errors.New("сумма перевода должна быть положительной")
	}
	return nil
}
