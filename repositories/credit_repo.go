package repositories

import (
	"bank-api/db"
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
	if err := c.DB.Create(credit).Error; err != nil {
		return err
	}
	return nil
}

func (c *creditRepository) GetCredits(userIDUint uint) ([]models.Credit, error) {
	var credits []models.Credit
	if err := db.DB.Where("user_id = ?", userIDUint).Find(&credits).Error; err != nil {
		return nil, err
	}
	return credits, nil
}
