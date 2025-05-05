package handlers

import (
	"bank-api/db"
	"bank-api/repositories"
	"bank-api/services"
	"encoding/json"
	"fmt"
	"net/http"
)

// TransferFunds — обработчик перевода средств
func TransferFunds(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var req struct {
		FromAccountID uint   `json:"from_account_id"`
		ToAccountID   uint   `json:"to_account_id"`
		Amount        int64  `json:"amount"`
		Description   string `json:"description,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FromAccountID == 0 || req.ToAccountID == 0 || req.Amount <= 0 {
		http.Error(w, "Invalid transfer parameters", http.StatusBadRequest)
		return
	}

	accRepo := repositories.GetAccountRepository(db.DB)
	transferService := services.NewTransactionService(db.DB, accRepo)

	err := transferService.Transfer(userID, req.FromAccountID, req.ToAccountID, req.Amount, req.Description)
	if err != nil {
		http.Error(w, fmt.Sprintf("Transfer failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transfer successful"})
}
