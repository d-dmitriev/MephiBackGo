package repositories

import (
	"bank-api/models"
	"gorm.io/gorm"
)

// cardRepository — реализация CardRepository через GORM
type cardRepository struct {
	DB *gorm.DB
}

var globalCardRepo CardRepository

// GetCardRepository — фабрика для создания репозитория
func GetCardRepository(db *gorm.DB) CardRepository {
	if globalCardRepo == nil {
		globalCardRepo = &cardRepository{DB: db}
	}
	return globalCardRepo
}

func (c *cardRepository) IssueCard(card *models.Card) (*models.Card, error) {
	if err := c.DB.Create(card).Error; err != nil {
		return nil, err
	}
	return card, nil
}

func (c *cardRepository) GetCards(userIDUint uint) ([]models.Card, error) {
	var cards []models.Card
	if err := c.DB.Where("user_id = ?", userIDUint).Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}
