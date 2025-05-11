package services

import (
	"Zira/internal/models"
	"Zira/internal/repository"
	"errors"

	"gorm.io/gorm"
)

// accountService реализует бизнес-логику для счетов
type accountService struct {
	repo repository.AccountRepository
}

// NewAccountService создает новый сервис счетов
func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

// CreateAccount создает новый счет для пользователя
func (s *accountService) CreateAccount(userID int64, accountType string) (interface{}, error) {
	acc := &models.Account{
		UserID:   uint(userID),
		Balance:  0,
		Currency: "RUB",
	}
	if err := s.repo.Create(acc); err != nil {
		return nil, err
	}
	return acc, nil
}

// GetUserAccounts возвращает все счета пользователя
func (s *accountService) GetUserAccounts(userID int64) (interface{}, error) {
	accounts, err := s.repo.GetByUserID(uint(userID))
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetAccountByID возвращает счет по ID
func (s *accountService) GetAccountByID(accountID int64) (interface{}, error) {
	acc, err := s.repo.GetByID(uint(accountID))
	if err != nil {
		return nil, err
	}
	return acc, nil
}

// Transfer реализует перевод между счетами
func (s *accountService) Transfer(req *models.TransferRequest) (interface{}, error) {
	// Получаем *gorm.DB из репозитория
	db, ok := getDBFromRepo(s.repo)
	if !ok {
		return nil, errors.New("Ошибка доступа к базе данных")
	}
	// Транзакция
	var result interface{}
	err := db.Transaction(func(tx *gorm.DB) error {
		from, err := s.repo.GetByID(req.FromAccountID)
		if err != nil {
			return errors.New("Счет отправителя не найден")
		}
		to, err := s.repo.GetByID(req.ToAccountID)
		if err != nil {
			return errors.New("Счет получателя не найден")
		}
		if from.Balance < req.Amount {
			return errors.New("Недостаточно средств")
		}
		from.Balance -= req.Amount
		to.Balance += req.Amount
		if err := tx.Save(from).Error; err != nil {
			return err
		}
		if err := tx.Save(to).Error; err != nil {
			return err
		}
		result = map[string]string{"status": "success"}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getDBFromRepo извлекает *gorm.DB из репозитория (временное решение)
func getDBFromRepo(repo interface{}) (*gorm.DB, bool) {
	if r, ok := repo.(interface{ DB() *gorm.DB }); ok {
		return r.DB(), true
	}
	return nil, false
}
