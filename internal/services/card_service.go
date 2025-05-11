package services

import (
	"Zira/internal/models"
	"Zira/internal/repository"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// CardService реализует бизнес-логику для банковских карт
type CardService struct {
	repo         repository.CardRepository
	pgpPublicKey string // путь или строка публичного PGP-ключа
	hmacSecret   []byte
}

// NewCardService создает новый сервис карт
func NewCardService(repo repository.CardRepository, pgpPublicKey string, hmacSecret []byte) *CardService {
	return &CardService{repo: repo, pgpPublicKey: pgpPublicKey, hmacSecret: hmacSecret}
}

// GenerateCard генерирует новую карту с валидным номером (алгоритм Луна), шифрует данные и хеширует CVV
func (s *CardService) GenerateCard(accountID uint, expiryMonth, expiryYear int, cvv string) (*models.Card, error) {
	number := generateLuhnCardNumber()
	cvvHash, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Ошибка хеширования CVV: %v", err)
	}
	// Здесь должна быть реализация PGP-шифрования номера и срока действия (заглушка)
	encrypted := fmt.Sprintf("ENCRYPTED(%s/%02d%02d)", number, expiryMonth, expiryYear)
	// Генерация HMAC для номера карты
	hmacValue := computeHMAC(number, s.hmacSecret)
	card := &models.Card{
		AccountID:    accountID,
		Number:       number, // В реальном проекте хранить только зашифрованно!
		CVVEncrypted: string(cvvHash),
		ExpiryMonth:  expiryMonth,
		ExpiryYear:   expiryYear,
		// Дополнительные поля для безопасности
		EncryptedData: encrypted,
		HMAC:          hmacValue,
	}
	if err := s.repo.Create(card); err != nil {
		return nil, err
	}
	return card, nil
}

// generateLuhnCardNumber генерирует валидный номер карты по алгоритму Луна
func generateLuhnCardNumber() string {
	rand.Seed(time.Now().UnixNano())
	num := make([]int, 16)
	for i := 0; i < 15; i++ {
		num[i] = rand.Intn(10)
	}
	// Алгоритм Луна для последней цифры
	sum := 0
	for i := 0; i < 15; i++ {
		digit := num[14-i]
		if i%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	num[15] = (10 - (sum % 10)) % 10
	result := ""
	for _, d := range num {
		result += fmt.Sprint(d)
	}
	return result
}

// computeHMAC вычисляет HMAC-SHA256 для строки
func computeHMAC(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
