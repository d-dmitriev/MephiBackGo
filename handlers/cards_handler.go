package handlers

import (
	"bank-api/db"
	"bank-api/repositories"
	"bank-api/services"
	"bank-api/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// IssueCardHandler — обработчик выпуска карты
func IssueCard(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value("userID").(string)
	accountIDStr := r.URL.Query().Get("account_id")

	if accountIDStr == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	userIDUint, err := utils.ParseUserID(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	acc := repositories.GetCardRepository(db.DB)
	cardService := services.NewCardService(acc)

	card, err := cardService.IssueCard(userIDUint, uint(accountID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to issue card: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}

func GetCards(w http.ResponseWriter, r *http.Request) {
	userIDInterface := r.Context().Value("userID")
	if userIDInterface == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDInterface.(string)
	if !ok {
		http.Error(w, "Invalid user ID type", http.StatusInternalServerError)
		return
	}

	card := repositories.GetCardRepository(db.DB)
	cardService := services.NewCardService(card)

	cards, err := cardService.GetCards(userID)

	if err != nil {
		http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cards)
}
