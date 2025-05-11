package models

import "time"

type Card struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	AccountID     uint   `gorm:"not null;index"`
	Number        string `gorm:"size:16;not null;uniqueIndex"` // PAN
	CVVEncrypted  string `gorm:"type:text;not null"`           // зашифрованный PGP CVV
	ExpiryMonth   int    `gorm:"not null"`
	ExpiryYear    int    `gorm:"not null"`
	EncryptedData string `gorm:"type:text" json:"encrypted_data"` // зашифрованные данные карты (PGP)
	HMAC          string `gorm:"type:text" json:"hmac"`           // HMAC для проверки целостности
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
