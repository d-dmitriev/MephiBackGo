package repositories

import (
	"bank-api/models"
	"gorm.io/gorm"
)

// cardRepository — реализация CardRepository через GORM
type creditRepository struct {
	DB *gorm.DB
}

var globalCreditRepo CreditRepository

// GetCreditRepository — фабрика для создания репозитория
func GetCreditRepository(db *gorm.DB) CreditRepository {
	if globalCreditRepo == nil {
		globalCreditRepo = &creditRepository{DB: db}
	}
	return globalCreditRepo
}

func (c *creditRepository) Create(credit *models.Credit) error {
	result := c.DB.Create(credit)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
