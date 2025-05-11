package repository

import (
	"Zira/internal/models"

	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(account *models.Account) error
	GetByID(id uint) (*models.Account, error)
	Update(account *models.Account) error
	GetByUserID(userID uint) ([]*models.Account, error)
}

type accountRepo struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepo{db: db}
}

func (r *accountRepo) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepo) GetByID(id uint) (*models.Account, error) {
	var a models.Account
	if err := r.db.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *accountRepo) Update(account *models.Account) error {
	return r.db.Save(account).Error
}

func (r *accountRepo) GetByUserID(userID uint) ([]*models.Account, error) {
	var accounts []*models.Account
	if err := r.db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
