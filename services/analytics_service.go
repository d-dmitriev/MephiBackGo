package services

import (
	"bank-api/utils"
	"fmt"
	"time"

	"bank-api/db"
	"bank-api/models"
)

type AnalyticsService struct{}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{}
}

// GetMonthlyStats — получает статистику доходов и расходов за последний месяц
func (s *AnalyticsService) GetMonthlyStats(userID string) (*models.MonthlyStats, error) {
	userIDUint, err := utils.ParseUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var transactions []models.Transaction
	result := db.DB.Where("sender_user_id = ? OR receiver_user_id = ?", userIDUint, userIDUint).
		Where("created_at >= ?", startOfMonth.Format(time.RFC3339)).
		Find(&transactions)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get transactions: %v", result.Error)
	}

	stats := &models.MonthlyStats{
		Month:         now.Format("January 2006"),
		TotalIncome:   0,
		TotalExpenses: 0,
		Transactions:  make([]models.TransactionSummary, 0),
	}

	for _, t := range transactions {
		summary := models.TransactionSummary{
			ID:          t.ID,
			Amount:      t.Amount,
			Type:        t.Type,
			Description: t.Description,
			Date:        t.CreatedAt,
		}

		if t.ReceiverUserID == userIDUint {
			stats.TotalIncome += t.Amount
			summary.Direction = "in"
		} else if t.SenderUserID == userIDUint {
			stats.TotalExpenses += t.Amount
			summary.Direction = "out"
		}

		stats.Transactions = append(stats.Transactions, summary)
	}

	stats.Balance = stats.TotalIncome - stats.TotalExpenses
	return stats, nil
}

// GetCreditLoad — получает информацию о кредитной нагрузке пользователя
func (s *AnalyticsService) GetCreditLoad(userID string) (*models.CreditLoad, error) {
	userIDUint, err := utils.ParseUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	var credits []models.Credit
	result := db.DB.Where("user_id = ?", userIDUint).Find(&credits)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get credits: %v", result.Error)
	}

	load := &models.CreditLoad{
		TotalCredits:      0,
		OutstandingDebt:   0,
		UpcomingPayments:  0,
		CreditsInProgress: 0,
		AverageInterest:   0,
		CreditUtilization: 0,
	}

	now := time.Now().Format(time.RFC3339)

	var totalInterest float64
	var activeCount int

	for _, c := range credits {
		load.TotalCredits += c.Amount
		load.OutstandingDebt += c.Amount - c.PaidAmount

		if c.Status == "active" {
			load.CreditsInProgress++
			totalInterest += c.InterestRate
			activeCount++
		}
	}

	if activeCount > 0 {
		load.AverageInterest = totalInterest / float64(activeCount)
	}

	if load.TotalCredits > 0 {
		load.CreditUtilization = float64(load.OutstandingDebt) / float64(load.TotalCredits) * 100
	}

	// Получаем ближайшие платежи
	var schedules []models.PaymentSchedule
	db.DB.Where("credit_id IN (?)", db.DB.Model(&models.Credit{}).Where("user_id = ?", userIDUint).Select("id")).
		Where("status = 'pending' AND due_date <= ?", now).
		Order("due_date ASC").
		Limit(5).
		Find(&schedules)

	for _, s := range schedules {
		load.UpcomingPayments += s.AmountDue
	}

	return load, nil
}

// PredictBalance — прогнозирует баланс пользователя на N дней вперед
func (s *AnalyticsService) PredictBalance(userID string, days int) (*models.BalancePrediction, error) {
	if days < 1 || days > 365 {
		return nil, fmt.Errorf("prediction period must be between 1 and 365 days")
	}

	userIDUint, err := utils.ParseUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	var accounts []models.Account
	db.DB.Where("user_id = ?", userIDUint).Find(&accounts)

	prediction := &models.BalancePrediction{
		Days:           days,
		StartDate:      time.Now().Format(time.RFC3339),
		EndDate:        time.Now().AddDate(0, 0, days).Format(time.RFC3339),
		BalanceHistory: make([]models.BalancePoint, 0),
	}

	for _, acc := range accounts {
		history := models.BalancePoint{
			AccountID: acc.ID,
			History:   make([]models.BalanceEntry, 0),
		}

		balance := acc.Balance
		history.History = append(history.History, models.BalanceEntry{
			Date:    time.Now().Format("2006-01-02"),
			Balance: balance,
		})

		for i := 1; i <= days; i++ {
			date := time.Now().AddDate(0, 0, i)
			dayOfWeek := date.Weekday()

			// Пример: каждый понедельник +500р, среду -300р, пятницу -1000р
			switch dayOfWeek {
			case time.Monday:
				balance += 50000 // Зарплата
			case time.Wednesday:
				balance -= 30000 // Коммунальные услуги
			case time.Friday:
				balance -= 100000 // Покупки
			}

			history.History = append(history.History, models.BalanceEntry{
				Date:    date.Format("2006-01-02"),
				Balance: balance,
			})
		}

		prediction.BalanceHistory = append(prediction.BalanceHistory, history)
	}

	return prediction, nil
}
