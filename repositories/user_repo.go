package repositories

import (
	"bank-api/models"
	"gorm.io/gorm"
)

// userRepository — реализация UserRepository через GORM
type userRepository struct {
	DB *gorm.DB
}

var globalUserRepo UserRepository

// GetUserRepository — фабрика для создания репозитория
func GetUserRepository(db *gorm.DB) UserRepository {
	if globalUserRepo == nil {
		globalUserRepo = &userRepository{DB: db}
	}
	return globalUserRepo
}

// Create — сохраняет нового пользователя в БД
func (u *userRepository) Create(user *models.User) error {
	result := u.DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindByEmail — ищет пользователя по email
func (u *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := u.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByID — ищет пользователя по ID
func (u *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	result := u.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
