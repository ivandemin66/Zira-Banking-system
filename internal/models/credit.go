package models

import "time"

type Credit struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	AccountID    uint    `gorm:"not null;index"`
	Amount       float64 `gorm:"type:numeric;not null"`
	InterestRate float64 `gorm:"type:numeric;not null"` // годовая процентная ставка
	TermMonths   int     `gorm:"not null"`
	Status       string  `gorm:"size:20;not null"` // "active", "closed" и т.п.
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Payments []PaymentSchedule `gorm:"foreignKey:CreditID"`
}
