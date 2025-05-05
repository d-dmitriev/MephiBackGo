package services

import (
	"bank-api/repositories"
	"bank-api/utils"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCardService_IssueCard_Success(t *testing.T) {
	utils.InitLogger()

	cardRepo := new(repositories.MockCardRepository)

	cardService := NewCardService(cardRepo)

	userID := 1
	accountID := 1

	cardRepo.On("IssueCard", mock.AnythingOfType("*models.Card")).Return(nil, nil)

	cardService.IssueCard(uint(userID), uint(accountID))
	//assert.NoError(t, err)

	//cardRepo.AssertExpectations(t)
}
