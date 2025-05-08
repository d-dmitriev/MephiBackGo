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
	if err := u.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// FindByEmail — ищет пользователя по email
func (u *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID — ищет пользователя по ID
func (u *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := u.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
