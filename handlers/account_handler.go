package handlers

import (
	"bank-api/db"
	"bank-api/repositories"
	"bank-api/services"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// CreateAccountHandler — обработчик создания банковского счёта
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста безопасно
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

	// Получаем тип счёта из JSON-тела
	var req struct {
		Type string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Type == "" {
		http.Error(w, "Missing account type", http.StatusBadRequest)
		return
	}

	acc := repositories.GetAccountRepository(db.DB)
	txn := repositories.GetTransactionRepository(db.DB)
	accountService := services.NewAccountService(acc, txn)

	err := accountService.CreateAccount(userID, req.Type)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create account: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Account created"})
}

func GetAccounts(w http.ResponseWriter, r *http.Request) {
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

	acc := repositories.GetAccountRepository(db.DB)
	txn := repositories.GetTransactionRepository(db.DB)
	accountService := services.NewAccountService(acc, txn)

	accounts, err := accountService.GetAccounts(userID)

	if err != nil {
		http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accounts)
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	// Получаем ID счёта из URL
	vars := mux.Vars(r)
	accountIDStr := vars["id"]

	// Конвертируем в uint
	accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	// Получаем тип счёта из JSON-тела
	var req struct {
		Amount int64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	acc := repositories.GetAccountRepository(db.DB)
	txn := repositories.GetTransactionRepository(db.DB)
	accountService := services.NewAccountService(acc, txn)

	if req.Amount > 0 {
		if err := accountService.Deposit(uint(accountID), req.Amount); err != nil {
			http.Error(w, "Failed to update balance: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if req.Amount < 0 {
		if err := accountService.Withdraw(uint(accountID), -req.Amount); err != nil {
			http.Error(w, "Failed to update balance: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":        "Balance updated successfully",
		"balance_change": strconv.FormatInt(req.Amount, 10),
	})
}
