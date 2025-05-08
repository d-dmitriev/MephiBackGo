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

// NewCreditRepository — фабрика для создания репозитория
func NewCreditRepository(db *gorm.DB) CreditRepository {
	return &creditRepository{DB: db}
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
