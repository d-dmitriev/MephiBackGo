package repositories

import (
	"bank-api/models"
	"fmt"
	"gorm.io/gorm"
)

// paymentScheduleRepository — реализация через GORM
type paymentScheduleRepository struct {
	DB *gorm.DB
}

var globalPaymentScheduleRepo PaymentScheduleRepository

func GetPaymentScheduleRepository(db *gorm.DB) PaymentScheduleRepository {
	if globalPaymentScheduleRepo == nil {
		globalPaymentScheduleRepo = &paymentScheduleRepository{DB: db}
	}
	return globalPaymentScheduleRepo
}

// Create — сохраняет новый график платежа
func (r *paymentScheduleRepository) Create(schedule *models.PaymentSchedule) error {
	result := r.DB.Create(schedule)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByCreditID — получает все записи графика по ID кредита
func (r *paymentScheduleRepository) GetByCreditID(creditID uint) ([]models.PaymentSchedule, error) {
	var schedules []models.PaymentSchedule
	result := r.DB.Where("credit_id = ?", creditID).Find(&schedules)
	if result.Error != nil {
		return nil, result.Error
	}
	return schedules, nil
}

// GetPendingByUserID — получает все незавершённые платежи пользователя
func (r *paymentScheduleRepository) GetPendingByUserID(userID uint) ([]models.PaymentSchedule, error) {
	var schedules []models.PaymentSchedule
	result := r.DB.Table("payment_schedules").
		Joins("JOIN credits ON payment_schedules.credit_id = credits.id").
		Where("credits.user_id = ? AND payment_schedules.status IN ('pending', 'overdue')", userID).
		Order("due_date ASC").
		Find(&schedules)

	if result.Error != nil {
		return nil, result.Error
	}
	return schedules, nil
}

// Update — обновляет запись в графике платежей
func (r *paymentScheduleRepository) Update(schedule *models.PaymentSchedule) error {
	result := r.DB.Save(schedule)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetPendingPayments — получает все платежи со статусом "pending" или "overdue"
func (r *paymentScheduleRepository) GetPendingPayments() ([]models.PaymentSchedule, error) {
	var schedules []models.PaymentSchedule

	// Выбираем только те платежи, у которых статус pending или overdue
	result := r.DB.Where("status IN (?)", []string{"pending", "overdue"}).Find(&schedules)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get pending payments: %v", result.Error)
	}

	return schedules, nil
}
