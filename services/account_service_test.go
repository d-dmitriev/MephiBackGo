// services/account_service_test.go

package services

import (
	"bank-api/repositories"
	"bank-api/utils"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestCreateAccount_Success — успешное создание счёта
func TestCreateAccount_Success(t *testing.T) {
	utils.InitLogger()

	accRepo := new(repositories.MockAccountRepository)
	txnRepo := new(repositories.MockTransactionRepository)

	accountService := NewAccountService(accRepo, txnRepo)

	userID := "1"
	accountType := "debit"

	// Настроить поведение мока
	accRepo.On("Create", mock.AnythingOfType("*models.Account")).Return(nil)

	err := accountService.CreateAccount(userID, accountType)

	assert.NoError(t, err)

	accRepo.AssertExpectations(t)
}

// TestCreateAccount_InvalidUserID — невалидный userID
func TestCreateAccount_InvalidUserID(t *testing.T) {
	utils.InitLogger()

	accRepo := new(repositories.MockAccountRepository)

	accountService := &AccountService{
		accountRepo: accRepo,
	}

	invalidUserID := "invalid"
	accountType := "debit"

	err := accountService.CreateAccount(invalidUserID, accountType)

	assert.Error(t, err)
}

// TestCreateAccount_AccountCreationFailed — ошибка при сохранении счёта
func TestCreateAccount_AccountCreationFailed(t *testing.T) {
	utils.InitLogger()

	accRepo := new(repositories.MockAccountRepository)

	accountService := &AccountService{
		accountRepo: accRepo,
	}

	userID := "1"
	accountType := "debit"

	// Настроить мок так, чтобы Create вернул ошибку
	accRepo.On("Create", mock.AnythingOfType("*models.Account")).Return(errors.New("DB error"))

	err := accountService.CreateAccount(userID, accountType)

	assert.Error(t, err)

	accRepo.AssertExpectations(t)
}
