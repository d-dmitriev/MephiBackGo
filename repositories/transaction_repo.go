package repositories

import (
	"bank-api/models"
	"gorm.io/gorm"
)

// transactionRepository — реализация TransactionRepository через GORM
type transactionRepository struct {
	DB *gorm.DB
}

var globalTransactionRepo TransactionRepository

// GetTransactionRepository — фабрика для создания репозитория
func GetTransactionRepository(db *gorm.DB) TransactionRepository {
	if globalTransactionRepo == nil {
		globalTransactionRepo = &transactionRepository{DB: db}
	}
	return globalTransactionRepo
}

func (t *transactionRepository) Create(transaction *models.Transaction) error {
	result := t.DB.Create(transaction)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (t *transactionRepository) FindByID(id uint) (*models.Transaction, error) {
	var acc models.Transaction
	if err := t.DB.First(&acc, id).Error; err != nil {
		return nil, err
	}
	return &acc, nil
}

func (t *transactionRepository) GetAllByUserID(userID uint) ([]models.Transaction, error) {
	var accounts []models.Transaction
	result := t.DB.Where("user_id = ?", userID).Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}
	return accounts, nil
}
