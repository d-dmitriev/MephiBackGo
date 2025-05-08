package services

import (
	"bank-api/integrations"
	"bank-api/repositories"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"

	"bank-api/models"
	"bank-api/utils"
)

type TransactionService struct {
	DB          *gorm.DB
	accountRepo repositories.AccountRepository
	userRepo    repositories.UserRepository
}

func NewTransactionService(
	db *gorm.DB,
	accountRepo repositories.AccountRepository,
	userRepo repositories.UserRepository,
) *TransactionService {
	return &TransactionService{
		DB:          db,
		accountRepo: accountRepo,
		userRepo:    userRepo,
	}
}

// Transfer — перевод средств между счетами
func (t *TransactionService) Transfer(userID string, fromAccountID, toAccountID uint, amount int64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	userIDUint, err := utils.ParseUserID(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	fromAccount, err := t.accountRepo.GetByID(fromAccountID)
	if err != nil {
		return fmt.Errorf("from account not found: %v", err)
	}
	toAccount, err := t.accountRepo.GetByID(toAccountID)
	if err != nil {
		return fmt.Errorf("to account not found: %v", err)
	}

	// Проверяем, что пользователь владеет счётом отправителя
	if fromAccount.UserID != userIDUint {
		return errors.New("access denied to sender account")
	}

	// Проверяем баланс
	if fromAccount.Balance < amount {
		return errors.New("insufficient funds")
	}

	// Начинаем транзакцию
	tx := t.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Списание со счёта отправителя
	fromAccount.Balance -= amount
	if err := fromAccount.UpdateBalance(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to deduct from sender: %v", err)
	}

	// Пополнение счёта получателя
	toAccount.Balance += amount
	if err := toAccount.UpdateBalance(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add to receiver: %v", err)
	}

	// Логируем транзакцию
	transaction := &models.Transaction{
		FromAccountID:  fromAccountID,
		ToAccountID:    toAccountID,
		SenderUserID:   userIDUint,
		ReceiverUserID: toAccount.UserID,
		Amount:         amount,
		Type:           "transfer",
		Description:    description,
		Status:         "success",
		CreatedAt:      time.Now().Format(time.RFC3339),
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to log transaction: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Отправляем email-уведомление
	user, _ := t.userRepo.FindByID(userIDUint)
	if user != nil {
		integrations.SendEmail(
			user.Email,
			"Перевод успешно выполнен",
			fmt.Sprintf(`
                <h2>Успешный перевод</h2>
                <p>Со счёта %d на счёт %d переведено %.2f руб.</p>
                <small>%s</small>
            `, fromAccountID, toAccountID, float64(amount)/100, description),
		)
	}

	return nil
}
