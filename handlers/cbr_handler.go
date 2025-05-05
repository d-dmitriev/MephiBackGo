package handlers

import (
	"bank-api/db"
	"bank-api/repositories"
	"encoding/json"
	"fmt"
	"net/http"

	"bank-api/services"
)

func GetKeyRate(w http.ResponseWriter, r *http.Request) {
	accRepo := repositories.GetAccountRepository(db.DB)
	paymentRepo := repositories.GetPaymentScheduleRepository(db.DB)
	creditRepo := repositories.GetCreditRepository(db.DB)
	creditService := services.NewCreditService(paymentRepo, accRepo, creditRepo)
	rate, err := creditService.GetKeyRateWithMargin()
	if err != nil {
		http.Error(w, "Failed to get key rate: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"key_rate": fmt.Sprintf("%.2f%%", rate),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
