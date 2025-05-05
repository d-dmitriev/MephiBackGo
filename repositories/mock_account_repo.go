package repositories

import (
	"bank-api/models"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(account *models.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountRepository) GetByID(id uint) (*models.Account, error) {
	args := m.Called(id)
	if account, ok := args.Get(0).(*models.Account); ok {
		return account, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAccountRepository) Update(account *models.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountRepository) GetAccounts(userID uint) ([]models.Account, error) {
	args := m.Called(userID)
	if accounts, ok := args.Get(0).([]models.Account); ok {
		return accounts, args.Error(1)
	}
	return nil, args.Error(1)
}
