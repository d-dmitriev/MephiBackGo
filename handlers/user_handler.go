package handlers

import (
	"bank-api/db"
	"bank-api/repositories"
	"bank-api/utils"
	"encoding/json"
	"net/http"

	"bank-api/services"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.Logger.Warnf("Invalid request: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	repo := repositories.GetUserRepository(db.DB)
	authService := services.NewAuthService(repo)
	err := authService.Register(input.Email, input.Username, input.Password)
	if err != nil {
		utils.Logger.Errorf("Registration failed: %v", err)
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	utils.Logger.Infof("User registered: %s", input.Email)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	repo := repositories.GetUserRepository(db.DB)
	authService := services.NewAuthService(repo)
	token, err := authService.Login(input.Email, input.Password)
	if err != nil {
		http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
