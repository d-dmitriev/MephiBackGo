package repositories

import (
	"bank-api/db"
	"bank-api/models"
	"gorm.io/gorm"
	"time"
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
	var transaction models.Transaction
	if err := t.DB.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (t *transactionRepository) GetAllByUserID(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := t.DB.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (t *transactionRepository) GetAllBySenderOrReceiverAfterDate(userIDUint uint, date time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := db.DB.Where("sender_user_id = ? OR receiver_user_id = ?", userIDUint, userIDUint).
		Where("created_at >= ?", date.Format(time.RFC3339)).
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
