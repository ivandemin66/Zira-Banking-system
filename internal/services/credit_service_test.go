package services

import (
	"Zira/internal/models"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Моки для зависимостей CreditService

type mockCreditRepo struct {
	credits            []*models.Credit
	createErr          error
	getAllActiveResult []*models.Credit
	getAllActiveErr    error
}

func (m *mockCreditRepo) Create(credit *models.Credit) error {
	if m.createErr != nil {
		return m.createErr
	}
	credit.ID = uint(len(m.credits) + 1)
	m.credits = append(m.credits, credit)
	return nil
}
func (m *mockCreditRepo) GetByID(id uint) (*models.Credit, error) { return nil, nil }
func (m *mockCreditRepo) Update(credit *models.Credit) error      { return nil }
func (m *mockCreditRepo) GetAllActive() ([]*models.Credit, error) {
	if m.getAllActiveErr != nil {
		return nil, m.getAllActiveErr
	}
	return m.getAllActiveResult, nil
}

type mockScheduleRepo struct {
	schedules           []*models.PaymentSchedule
	createErr           error
	getByCreditIDResult []*models.PaymentSchedule
	getByCreditIDErr    error
	updateCalled        []*models.PaymentSchedule
}

func (m *mockScheduleRepo) Create(s *models.PaymentSchedule) error {
	if m.createErr != nil {
		return m.createErr
	}
	s.ID = uint(len(m.schedules) + 1)
	m.schedules = append(m.schedules, s)
	return nil
}
func (m *mockScheduleRepo) GetByCreditID(creditID uint) ([]*models.PaymentSchedule, error) {
	if m.getByCreditIDErr != nil {
		return nil, m.getByCreditIDErr
	}
	return m.getByCreditIDResult, nil
}
func (m *mockScheduleRepo) Update(s *models.PaymentSchedule) error {
	m.updateCalled = append(m.updateCalled, s)
	return nil
}

type mockAccountRepo struct {
	accounts map[uint]*models.Account
}

func (m *mockAccountRepo) GetByID(id uint) (*models.Account, error) {
	if acc, ok := m.accounts[id]; ok {
		return acc, nil
	}
	return nil, errors.New("not found")
}
func (m *mockAccountRepo) Create(a *models.Account) error                     { return nil }
func (m *mockAccountRepo) Update(a *models.Account) error                     { m.accounts[a.ID] = a; return nil }
func (m *mockAccountRepo) GetByUserID(userID uint) ([]*models.Account, error) { return nil, nil }

type mockEmailService struct {
	sent []string
}

func (m *mockEmailService) Send(to, subject, body string) error {
	m.sent = append(m.sent, subject+": "+body)
	return nil
}

// Тест: успешное оформление кредита
func TestCreditService_ApplyForCredit_Success(t *testing.T) {
	creditRepo := &mockCreditRepo{}
	scheduleRepo := &mockScheduleRepo{}
	accountRepo := &mockAccountRepo{}
	emailService := &mockEmailService{}
	service := NewCreditService(creditRepo, scheduleRepo, accountRepo, emailService)

	credit, err := service.ApplyForCredit(1, 100000, 12, 12)
	assert.NoError(t, err)
	assert.NotNil(t, credit)
	assert.Equal(t, 1, len(creditRepo.credits))
	assert.Equal(t, 12, len(scheduleRepo.schedules))
}

// Тест: ошибка при некорректных параметрах
func TestCreditService_ApplyForCredit_InvalidParams(t *testing.T) {
	service := NewCreditService(&mockCreditRepo{}, &mockScheduleRepo{}, &mockAccountRepo{}, nil)
	_, err := service.ApplyForCredit(1, 0, 12, 12)
	assert.Error(t, err)
	_, err = service.ApplyForCredit(1, 100000, 0, 12)
	assert.Error(t, err)
	_, err = service.ApplyForCredit(1, 100000, 12, 0)
	assert.Error(t, err)
}

// Тест: расчет аннуитетного платежа
func TestCalcAnnuity(t *testing.T) {
	monthly := calcAnnuity(100000, 12, 12)
	assert.InDelta(t, 8884.0, monthly, 1.0) // Проверяем, что расчет близок к ожидаемому
}

// Тест: обработка просроченных платежей (шедулер)
func TestCreditService_ProcessDuePayments(t *testing.T) {
	credit := &models.Credit{ID: 1, AccountID: 1, Amount: 100000, InterestRate: 12, TermMonths: 12, Status: "active"}
	acc := &models.Account{ID: 1, Balance: 9000, Currency: "RUB"}
	payment := &models.PaymentSchedule{ID: 1, CreditID: 1, DueDate: time.Now().AddDate(0, 0, -1), Amount: 8000, Paid: false}
	creditRepo := &mockCreditRepo{getAllActiveResult: []*models.Credit{credit}}
	scheduleRepo := &mockScheduleRepo{getByCreditIDResult: []*models.PaymentSchedule{payment}}
	accountRepo := &mockAccountRepo{accounts: map[uint]*models.Account{1: acc}}
	emailService := &mockEmailService{}
	service := NewCreditService(creditRepo, scheduleRepo, accountRepo, emailService)

	err := service.ProcessDuePayments()
	assert.NoError(t, err)
	assert.True(t, payment.Paid)
	assert.True(t, acc.Balance < 9000)
	assert.NotEmpty(t, emailService.sent)
}
