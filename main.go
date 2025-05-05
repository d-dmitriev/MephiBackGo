package main

import (
	"bank-api/config"
	"bank-api/db"
	"bank-api/models"
	"bank-api/utils"
	"log"
	"net/http"

	"bank-api/handlers"
	"bank-api/middleware"
	"github.com/gorilla/mux"
)

func main() {
	utils.InitLogger()
	utils.Logger.Info("Starting server...")

	cfg := config.LoadConfig()

	if err := db.Connect(cfg); err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	// Автоматическая миграция
	db.DB.AutoMigrate(&models.User{}, &models.Account{}, &models.Card{}, &models.Transaction{}, &models.Credit{}, &models.PaymentSchedule{})

	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/cbr/keyrate", handlers.GetKeyRate).Methods("GET")

	// Protected routes
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/accounts", handlers.GetAccounts).Methods("GET")
	protected.HandleFunc("/accounts", handlers.CreateAccount).Methods("POST")
	protected.HandleFunc("/accounts/{id}", handlers.UpdateAccount).Methods("PUT")
	protected.HandleFunc("/cards", handlers.GetCards).Methods("GET")
	protected.HandleFunc("/cards", handlers.IssueCard).Methods("POST")
	protected.HandleFunc("/transfer", handlers.TransferFunds).Methods("POST")
	protected.HandleFunc("/analytics", handlers.GetAnalytics).Methods("GET")
	protected.HandleFunc("/credits", handlers.ApplyForCredit).Methods("POST")

	finalRouter := middleware.LoggingMiddleware(r)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", finalRouter))
}
