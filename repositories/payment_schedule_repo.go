package repositories

import (
	"bank-api/db"
	"bank-api/models"
	"gorm.io/gorm"
	"time"
)

// paymentScheduleRepository — реализация через GORM
type paymentScheduleRepository struct {
	DB *gorm.DB
}

func NewPaymentScheduleRepository(db *gorm.DB) PaymentScheduleRepository {
	return &paymentScheduleRepository{DB: db}
}

// Create — сохраняет новый график платежа
func (r *paymentScheduleRepository) Create(schedule *models.PaymentSchedule) error {
	if err := r.DB.Create(schedule).Error; err != nil {
		return err
	}
	return nil
}

// GetByCreditID — получает все записи графика по ID кредита
func (r *paymentScheduleRepository) GetByCreditID(creditID uint) ([]models.PaymentSchedule, error) {
	var schedules []models.PaymentSchedule
	if err := r.DB.Where("credit_id = ?", creditID).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// GetPendingByUserID — получает все незавершённые платежи пользователя
func (r *paymentScheduleRepository) GetPendingByUserID(userID uint) ([]models.PaymentSchedule, error) {
	var schedules []models.PaymentSchedule
	if err := r.DB.Table("payment_schedules").
		Joins("JOIN credits ON payment_schedules.credit_id = credits.id").
		Where("credits.user_id = ? AND payment_schedules.status IN ('pending', 'overdue')", userID).
		Order("due_date ASC").
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// Update — обновляет запись в графике платежей
func (r *paymentScheduleRepository) Update(schedule *models.PaymentSchedule) error {
	if err := r.DB.Save(schedule).Error; err != nil {
		return err
	}
	return nil
}

// GetPendingPayments — получает все платежи со статусом "pending" или "overdue"
func (r *paymentScheduleRepository) GetPendingPayments() ([]models.PaymentSchedule, error) {
	var schedules []models.PaymentSchedule
	// Выбираем только те платежи, у которых статус pending или overdue
	if err := r.DB.Where("status IN (?)", []string{"pending", "overdue"}).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *paymentScheduleRepository) GetPendingPaymentsByUserAndDate(userIDUint uint, date time.Time) ([]models.PaymentSchedule, error) {
	var schedules []models.PaymentSchedule
	if err := db.DB.Where("credit_id IN (?)", db.DB.Model(&models.Credit{}).Where("user_id = ?", userIDUint).Select("id")).
		Where("status = 'pending' AND due_date <= ?", date.Format(time.RFC3339)).
		Order("due_date ASC").
		Limit(5).
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}
