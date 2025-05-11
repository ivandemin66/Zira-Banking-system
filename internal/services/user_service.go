package services

import (
	"Zira/internal/models"
	"Zira/internal/repository"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// UserService реализует бизнес-логику для пользователей
type userService struct {
	repo      repository.UserRepository
	jwtSecret string
}

// NewUserService создает новый сервис пользователей
func NewUserService(repo repository.UserRepository, jwtSecret string) UserService {
	return &userService{repo: repo, jwtSecret: jwtSecret}
}

// Register реализует регистрацию пользователя
func (s *userService) Register(req *models.RegisterRequest) (interface{}, error) {
	// Проверка уникальности email
	existing, _ := s.repo.GetByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("Пользователь с таким email уже существует")
	}
	// Хеширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Ошибка хеширования пароля: %v", err)
	}
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hash),
	}
	if err := s.repo.Create(user); err != nil {
		// Обработка ошибки уникальности email (Postgres 23505)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, fmt.Errorf("Email уже зарегистрирован (409)")
		}
		return nil, err
	}
	return user, nil
}

// Login реализует аутентификацию пользователя
func (s *userService) Login(req *models.LoginRequest) (interface{}, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("Пользователь не найден")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("Неверный пароль")
	}
	// Генерация JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   fmt.Sprint(user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("Ошибка генерации токена: %v", err)
	}
	return map[string]string{"token": tokenStr}, nil
}

// GetUserByID возвращает пользователя по ID
func (s *userService) GetUserByID(id int64) (interface{}, error) {
	user, err := s.repo.GetByID(uint(id))
	if err != nil {
		return nil, errors.New("Пользователь не найден")
	}
	return user, nil
}
