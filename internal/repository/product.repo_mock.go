package repository

import (
	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func (m *ProductRepositoryMock) CreateProduct(data *models.Product) (*models.Product, error) {
	args := m.Mock.Called()
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetAllProduct(query *models.ProductQuery) (*models.Products, int, error) {
	args := m.Mock.Called()
	return args.Get(0).(*models.Products), args.Get(1).(int), args.Error(2)
}

func (m *UserRepositoryMock) GetDetailProduct(uuid string) (*models.Product, error) {
	args := m.Mock.Called()
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *ProductRepositoryMock) UpdateProduct(uuid string, data *models.Product) (*models.Product, error) {
	args := m.Mock.Called()
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *ProductRepositoryMock) DeleteProduct(uuid string) (*models.Product, error) {
	args := m.Mock.Called()
	return args.Get(0).(*models.Product), args.Error(1)
}
