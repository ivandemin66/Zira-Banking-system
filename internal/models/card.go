package models

import "time"

type Card struct {
	ID          int64     `json:"id" db:"id"`
	AccountID   int64     `json:"account_id" db:"account_id"`
	Number      string    `json:"number" db:"number_encrypted"`
	ExpiryMonth int       `json:"expiry_month" db:"expiry_month_encrypted"`
	ExpiryYear  int       `json:"expiry_year" db:"expiry_year_encrypted"`
	CVVHash     string    `json:"-" db:"cvv_hash"`
	NumberHMAC  string    `json:"-" db:"number_hmac"`
	ExpiryHMAC  string    `json:"-" db:"expiry_hmac"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateCardRequest struct {
	AccountID int64 `json:"account_id"`
}
