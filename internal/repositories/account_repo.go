package repositories

import (
	"database/sql"
	"errors"
	"github.com/ivandemin66/Zira-Banking-system/internal/models"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(account *models.Account) (int64, error) {
	query := `
		INSERT INTO accounts (user_id, number, balance, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`
	var id int64
	err := r.db.QueryRow(
		query,
		account.UserID,
		account.Number,
		account.Balance,
		account.Type,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AccountRepository) GetByID(id int64) (*models.Account, error) {
	query := `
		SELECT id, user_id, number, balance, type, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`
	var account models.Account
	err := r.db.QueryRow(query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.Number,
		&account.Balance,
		&account.Type,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("счет не найден")
		}
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) GetByUserID(userID int64) ([]*models.Account, error) {
	query := `
		SELECT id, user_id, number, balance, type, created_at, updated_at
		FROM accounts
		WHERE user_id = $1
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*models.Account
	for rows.Next() {
		var account models.Account
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Number,
			&account.Balance,
			&account.Type,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	return accounts, nil
}

func (r *AccountRepository) UpdateBalance(id int64, amount float64, tx *sql.Tx) error {
	var query string
	query = `
		UPDATE accounts
		SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2
		RETURNING balance
	`

	var db interface {
		QueryRow(query string, args ...interface{}) *sql.Row
	}

	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	var newBalance float64
	err := db.QueryRow(query, amount, id).Scan(&newBalance)
	if err != nil {
		return err
	}

	if newBalance < 0 {
		return errors.New("недостаточно средств на счете")
	}

	return nil
}
