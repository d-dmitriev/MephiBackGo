package repositories

import (
	"bank-api/models"
	"github.com/stretchr/testify/mock"
)

type MockCardRepository struct {
	mock.Mock
}

func (m *MockCardRepository) IssueCard(card *models.Card) (*models.Card, error) {
	args := m.Called(card)
	if card, ok := args.Get(0).(*models.Card); ok {
		return card, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCardRepository) GetCards(userID uint) ([]models.Card, error) {
	args := m.Called(userID)
	if cards, ok := args.Get(0).([]models.Card); ok {
		return cards, args.Error(1)
	}
	return nil, args.Error(1)
}
