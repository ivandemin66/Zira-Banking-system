package repository

import (
	"Zira/internal/models"

	"gorm.io/gorm"
)

type CreditRepository interface {
	Create(credit *models.Credit) error
	GetByID(id uint) (*models.Credit, error)
	Update(credit *models.Credit) error
	GetAllActive() ([]*models.Credit, error)
}

type creditRepo struct {
	db *gorm.DB
}

func NewCreditRepository(db *gorm.DB) CreditRepository {
	return &creditRepo{db: db}
}

func (r *creditRepo) Create(credit *models.Credit) error {
	return r.db.Create(credit).Error
}

func (r *creditRepo) GetByID(id uint) (*models.Credit, error) {
	var c models.Credit
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *creditRepo) Update(credit *models.Credit) error {
	return r.db.Save(credit).Error
}

func (r *creditRepo) GetAllActive() ([]*models.Credit, error) {
	var credits []*models.Credit
	if err := r.db.Where("status = ?", "active").Find(&credits).Error; err != nil {
		return nil, err
	}
	return credits, nil
}
