package handlers

import (
	"bank-api/db"
	"bank-api/repositories"
	"bank-api/services"
	"bank-api/utils"
	"encoding/json"
	"net/http"
)

func ApplyForCredit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	userIDUint, err := utils.ParseUserID(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		AccountID   uint    `json:"account_id"`
		Amount      int64   `json:"amount"` // в копейках
		Rate        float64 `json:"rate"`   // например, 12.5
		DurationDay int     `json:"duration_day"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	accRepo := repositories.GetAccountRepository(db.DB)
	paymentRepo := repositories.GetPaymentScheduleRepository(db.DB)
	creditRepo := repositories.GetCreditRepository(db.DB)
	creditService := services.NewCreditService(paymentRepo, accRepo, creditRepo)

	credit, err := creditService.IssueCredit(userIDUint, req.AccountID, req.Amount, req.Rate, req.DurationDay)
	if err != nil {
		http.Error(w, "Failed to issue credit: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(credit)
}
