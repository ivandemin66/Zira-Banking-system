package repository

import (
	"Zira/internal/models"

	"gorm.io/gorm"
)

type PaymentScheduleRepository interface {
	Create(schedule *models.PaymentSchedule) error
	GetByCreditID(creditID uint) ([]*models.PaymentSchedule, error)
	Update(schedule *models.PaymentSchedule) error
}

type paymentScheduleRepo struct {
	db *gorm.DB
}

func NewPaymentScheduleRepository(db *gorm.DB) PaymentScheduleRepository {
	return &paymentScheduleRepo{db: db}
}

func (r *paymentScheduleRepo) Create(schedule *models.PaymentSchedule) error {
	return r.db.Create(schedule).Error
}

func (r *paymentScheduleRepo) GetByCreditID(creditID uint) ([]*models.PaymentSchedule, error) {
	var schedules []*models.PaymentSchedule
	if err := r.db.Where("credit_id = ?", creditID).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *paymentScheduleRepo) Update(schedule *models.PaymentSchedule) error {
	return r.db.Save(schedule).Error
}
