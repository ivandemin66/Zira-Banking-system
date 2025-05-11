package services

import (
	"Zira/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	users map[string]*models.User
}

func (m *mockUserRepo) Create(user *models.User) error {
	m.users[user.Email] = user
	return nil
}
func (m *mockUserRepo) GetByEmail(email string) (*models.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, nil
}
func (m *mockUserRepo) GetByID(id uint) (*models.User, error) { return nil, nil }

func TestUserService_Register(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]*models.User)}
	service := NewUserService(repo, "secret")

	// Регистрация нового пользователя
	req := &models.RegisterRequest{Name: "Ivan", Email: "ivan@example.com", Password: "password123"}
	user, err := service.Register(req)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Попытка регистрации с тем же email
	_, err = service.Register(req)
	assert.Error(t, err)
}
