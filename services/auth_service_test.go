package services

import (
	"bank-api/repositories"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"bank-api/models"
	"bank-api/utils"
)

func TestLogin_Success(t *testing.T) {
	utils.InitLogger()

	repo := new(repositories.MockUserRepository)
	service := NewAuthService(repo)

	validEmail := "test@example.com"
	validPass := "password123"

	hashedPass, err := utils.HashPassword(validPass)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	user := models.User{
		ID:       1,
		Email:    validEmail,
		Password: hashedPass,
	}

	// Настроить поведение мока
	repo.On("FindByEmail", validEmail).Return(&user, nil)

	token, err := service.Login(validEmail, validPass)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	repo.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	utils.InitLogger()

	repo := new(repositories.MockUserRepository)
	service := NewAuthService(repo)

	repo.On("FindByEmail", "wrong@example.com").Return(nil, errors.New("not found"))

	token, err := service.Login("wrong@example.com", "wrongpass")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if token != "" {
		t.Errorf("Expected empty token on invalid credentials")
	}
	repo.AssertExpectations(t)
}
