package services

import (
	"Zira/internal/models"
	"Zira/internal/repository"
	"errors"
	"math"
	"time"
)

// EmailSender определяет интерфейс для отправки email (для мокирования)
type EmailSender interface {
	Send(to, subject, body string) error
}

// CreditService реализует бизнес-логику для кредитов
// Использует EmailSender вместо *EmailService
type CreditService struct {
	repo         repository.CreditRepository
	scheduleRepo repository.PaymentScheduleRepository
	accountRepo  repository.AccountRepository
	emailService EmailSender // Используем интерфейс
}

// NewCreditService создает новый сервис кредитов
func NewCreditService(repo repository.CreditRepository, scheduleRepo repository.PaymentScheduleRepository, accountRepo repository.AccountRepository, emailService EmailSender) *CreditService {
	return &CreditService{repo: repo, scheduleRepo: scheduleRepo, accountRepo: accountRepo, emailService: emailService}
}

// ApplyForCredit оформляет кредит, рассчитывает график платежей
func (s *CreditService) ApplyForCredit(accountID uint, amount float64, months int, rate float64) (*models.Credit, error) {
	if amount <= 0 || months <= 0 || rate <= 0 {
		return nil, errors.New("Некорректные параметры кредита")
	}
	credit := &models.Credit{
		AccountID:    accountID,
		Amount:       amount,
		InterestRate: rate,
		TermMonths:   months,
		Status:       "active",
	}
	if err := s.repo.Create(credit); err != nil {
		return nil, err
	}
	// Генерация графика платежей
	monthly := calcAnnuity(amount, rate, months)
	for i := 1; i <= months; i++ {
		ps := &models.PaymentSchedule{
			CreditID: credit.ID,
			DueDate:  time.Now().AddDate(0, i, 0),
			Amount:   monthly,
			Paid:     false,
		}
		if err := s.scheduleRepo.Create(ps); err != nil {
			return nil, err
		}
	}
	return credit, nil
}

// calcAnnuity рассчитывает аннуитетный платеж
func calcAnnuity(sum, rate float64, months int) float64 {
	monthlyRate := rate / 100 / 12
	return math.Round(sum*monthlyRate/(1-math.Pow(1+monthlyRate, float64(-months)))*100) / 100
}

// ProcessDuePayments списывает платежи по кредитам (шедулер)
func (s *CreditService) ProcessDuePayments() error {
	credits, err := s.repo.GetAllActive()
	if err != nil {
		return err
	}
	for _, credit := range credits {
		payments, err := s.scheduleRepo.GetByCreditID(credit.ID)
		if err != nil {
			continue
		}
		for _, payment := range payments {
			if !payment.Paid && payment.DueDate.Before(time.Now()) {
				acc, err := s.accountRepo.GetByID(credit.AccountID)
				if err != nil {
					continue
				}
				if acc.Balance >= payment.Amount {
					acc.Balance -= payment.Amount
					payment.Paid = true
					payment.PaidAt = ptrTime(time.Now())
					_ = s.accountRepo.Update(acc)
					_ = s.scheduleRepo.Update(payment)
					// Отправить email об успешном платеже (заглушка)
					_ = s.emailService.Send("user@example.com", "Платеж по кредиту", "Платеж успешно списан")
				} else {
					// Начислить штраф (+10%)
					payment.Amount = payment.Amount * 1.1
					_ = s.scheduleRepo.Update(payment)
					// Отправить email о просрочке (заглушка)
					_ = s.emailService.Send("user@example.com", "Просрочка платежа по кредиту", "На вашем счете недостаточно средств. Начислен штраф.")
				}
			}
		}
	}
	return nil
}

func ptrTime(t time.Time) *time.Time { return &t }
