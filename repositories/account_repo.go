package repositories

import (
	"bank-api/models"
	"gorm.io/gorm"
)

// accountRepository — реализация AccountRepository через GORM
type accountRepository struct {
	DB *gorm.DB
}

var globalAccountRepo AccountRepository

// GetAccountRepository — фабрика для создания репозитория
func GetAccountRepository(db *gorm.DB) AccountRepository {
	if globalAccountRepo == nil {
		globalAccountRepo = &accountRepository{DB: db}
	}
	return globalAccountRepo
}

func (a *accountRepository) Create(account *models.Account) error {
	result := a.DB.Create(account)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *accountRepository) GetByID(id uint) (*models.Account, error) {
	var acc models.Account
	if err := a.DB.First(&acc, id).Error; err != nil {
		return nil, err
	}
	return &acc, nil
}

func (a *accountRepository) Update(account *models.Account) error {
	result := a.DB.Save(account)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *accountRepository) GetAccounts(userIDUint uint) ([]models.Account, error) {
	var accounts []models.Account
	result := a.DB.Where("user_id = ?", userIDUint).Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}
	return accounts, nil
}
