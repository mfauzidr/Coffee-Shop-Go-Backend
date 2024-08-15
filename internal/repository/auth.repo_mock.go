package repository

import (
	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type AuthRepositoryMock struct {
	mock.Mock
}

func (m *AuthRepositoryMock) RegisterUser(data *models.Users) (string, error) {
	args := m.Mock.Called()
	return args.Get(0).(string), args.Error(1)
}

func (m *AuthRepositoryMock) GetByEmail(email string) (*models.UsersRes, error) {
	args := m.Mock.Called()
	return args.Get(0).(*models.UsersRes), args.Error(1)
}
