package models

import "time"

type PaymentSchedule struct {
	ID       uint       `gorm:"primaryKey;autoIncrement"`
	CreditID uint       `gorm:"not null;index"`
	DueDate  time.Time  `gorm:"not null"`
	Amount   float64    `gorm:"type:numeric;not null"`
	Paid     bool       `gorm:"not null;default:false"`
	PaidAt   *time.Time // nil, если ещё не оплачено
}
