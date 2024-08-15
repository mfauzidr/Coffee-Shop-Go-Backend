package repository

import (
	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) CreateUser(data *models.Users) (*models.Users, error) {
	args := m.Mock.Called()
	return args.Get(0).(*models.Users), args.Error(1)
}

func (m *UserRepositoryMock) GetAllUsers(query *models.UsersQuery) (*models.UsersRes, int, error) {

	args := m.Mock.Called()
	return args.Get(0).(*models.UsersRes), args.Get(1).(int), args.Error(2)
}

func (m *UserRepositoryMock) GetDetailsUser(uuid string) (*models.Users, error) {
	args := m.Mock.Called()
	return args.Get(0).(*models.Users), args.Error(1)
}

func (m *UserRepositoryMock) UpdateUser(uuid string, data *models.Users) (*models.Users, error) {
	args := m.Mock.Called()
	return args.Get(0).(*models.Users), args.Error(1)
}

func (m *UserRepositoryMock) DeleteUser(uuid string) (*models.Users, error) {
	args := m.Mock.Called()
	return args.Get(0).(*models.Users), args.Error(1)
}
