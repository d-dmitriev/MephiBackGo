package models

type MonthlyStats struct {
	Month         string               `json:"month"`
	TotalIncome   int64                `json:"total_income"`   // в копейках
	TotalExpenses int64                `json:"total_expenses"` // в копейках
	Balance       int64                `json:"balance"`        // в копейках
	Transactions  []TransactionSummary `json:"transactions"`
}

type TransactionSummary struct {
	ID          uint   `json:"id"`
	Amount      int64  `json:"amount"`    // в копейках
	Direction   string `json:"direction"` // in / out
	Type        string `json:"type"`      // transfer, credit_payment
	Description string `json:"description"`
	Date        string `json:"date"` // RFC3339
}

type CreditLoad struct {
	TotalCredits      int64   `json:"total_credits"`       // сумма всех кредитов
	OutstandingDebt   int64   `json:"outstanding_debt"`    // долг
	UpcomingPayments  int64   `json:"upcoming_payments"`   // ближайшие платежи
	CreditsInProgress int     `json:"credits_in_progress"` // активных кредитов
	AverageInterest   float64 `json:"average_interest"`    // средняя ставка
	CreditUtilization float64 `json:"credit_utilization"`  // процент от общего лимита
}

type BalancePrediction struct {
	Days           int            `json:"days"`
	StartDate      string         `json:"start_date"` // RFC3339
	EndDate        string         `json:"end_date"`   // RFC3339
	BalanceHistory []BalancePoint `json:"balance_history"`
}

type BalancePoint struct {
	AccountID uint           `json:"-"`
	History   []BalanceEntry `json:"history"`
}

type BalanceEntry struct {
	Date    string `json:"date"`    // YYYY-MM-DD
	Balance int64  `json:"balance"` // в копейках
}
