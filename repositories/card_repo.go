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
	result := c.DB.Create(card)
	if result.Error != nil {
		return nil, result.Error
	}
	return card, nil
}

func (c *cardRepository) GetCards(userIDUint uint) ([]models.Card, error) {
	var cards []models.Card
	result := c.DB.Where("user_id = ?", userIDUint).Find(&cards)
	if result.Error != nil {
		return nil, result.Error
	}
	return cards, nil
}
