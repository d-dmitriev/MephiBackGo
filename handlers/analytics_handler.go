package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bank-api/services"
)

func GetAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		var err error
		days, err = strconv.Atoi(daysStr)
		if err != nil || days < 1 || days > 365 {
			http.Error(w, "Invalid days parameter", http.StatusBadRequest)
			return
		}
	}

	analyticsService := services.NewAnalyticsService()

	stats, _ := analyticsService.GetMonthlyStats(userID)
	creditLoad, _ := analyticsService.GetCreditLoad(userID)
	prediction, _ := analyticsService.PredictBalance(userID, days)

	response := map[string]interface{}{
		"monthly_stats":      stats,
		"credit_load":        creditLoad,
		"balance_prediction": prediction,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
