package repositories

import (
	"bank-api/models"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(txn *models.Transaction) error {
	args := m.Called(txn)
	return args.Error(0)
}

func (t *MockTransactionRepository) FindByID(id uint) (*models.Transaction, error) {
	args := t.Called(id)
	if txn, ok := args.Get(0).(*models.Transaction); ok {
		return txn, args.Error(1)
	}
	return nil, args.Error(1)
}

func (t *MockTransactionRepository) GetAllByUserID(userID uint) ([]models.Transaction, error) {
	args := t.Called(userID)
	if accounts, ok := args.Get(0).([]models.Transaction); ok {
		return accounts, args.Error(1)
	}
	return nil, args.Error(1)
}
