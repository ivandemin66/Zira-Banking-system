package models

import "time"

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
	Transfer   TransactionType = "transfer"
	Payment    TransactionType = "payment"
)

type Transaction struct {
	ID            int64           `json:"id" db:"id"`
	FromAccountID *int64          `json:"from_account_id,omitempty" db:"from_account_id"`
	ToAccountID   *int64          `json:"to_account_id,omitempty" db:"to_account_id"`
	Amount        float64         `json:"amount" db:"amount"`
	Type          TransactionType `json:"type" db:"type"`
	Description   string          `json:"description" db:"description"`
	Status        string          `json:"status" db:"status"` // "ожидание", "завершено", "сбой"
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
}
