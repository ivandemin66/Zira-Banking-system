package repository

import (
	"Zira/internal/models"

	"gorm.io/gorm"
)

type CardRepository interface {
	Create(card *models.Card) error
	GetByID(id uint) (*models.Card, error)
	GetByNumber(number string) (*models.Card, error)
}

type cardRepo struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) CardRepository {
	return &cardRepo{db: db}
}

func (r *cardRepo) Create(card *models.Card) error {
	return r.db.Create(card).Error
}

func (r *cardRepo) GetByID(id uint) (*models.Card, error) {
	var c models.Card
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *cardRepo) GetByNumber(number string) (*models.Card, error) {
	var c models.Card
	if err := r.db.Where("number = ?", number).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
