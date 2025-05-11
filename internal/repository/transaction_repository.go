package repository

import (
	"Zira/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *models.Transaction) error
	ListByAccount(accountID uint) ([]models.Transaction, error)
	GetByID(id uint) (*models.Transaction, error)
}

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepo) ListByAccount(accountID uint) ([]models.Transaction, error) {
	var txs []models.Transaction
	if err := r.db.
		Where("from_account_id = ? OR to_account_id = ?", accountID, accountID).
		Order("created_at desc").
		Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *transactionRepo) GetByID(id uint) (*models.Transaction, error) {
	var t models.Transaction
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
