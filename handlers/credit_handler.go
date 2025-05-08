package handlers

import (
	"bank-api/services"
	"bank-api/utils"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CreditHandler struct {
	creditService *services.CreditService
	logger        *logrus.Logger
}

func NewCreditHandler(creditService *services.CreditService, logger *logrus.Logger) *CreditHandler {
	return &CreditHandler{
		creditService: creditService,
		logger:        logger,
	}
}

func (c *CreditHandler) ApplyForCredit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	userIDUint, err := utils.ParseUserID(userID)
	if err != nil {
		c.logger.Warnf("Invalid user ID: %v", userIDUint)
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
		c.logger.Warnf("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	credit, err := c.creditService.IssueCredit(userIDUint, req.AccountID, req.Amount, req.Rate, req.DurationDay)
	if err != nil {
		c.logger.Warnf("Failed to issue credit: " + err.Error())
		http.Error(w, "Failed to issue credit: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(credit)
}
