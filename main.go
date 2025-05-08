package main

import (
	"bank-api/config"
	"bank-api/db"
	"bank-api/models"
	"bank-api/repositories"
	"bank-api/services"
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

	userRepo := repositories.NewUserRepository(db.DB)
	cardRepo := repositories.NewCardRepository(db.DB)
	creditRepo := repositories.NewCreditRepository(db.DB)
	accountRepo := repositories.NewAccountRepository(db.DB)
	transactionRepo := repositories.NewTransactionRepository(db.DB)
	paymentScheduleRepo := repositories.NewPaymentScheduleRepository(db.DB)

	authService := services.NewAuthService(userRepo)
	accountService := services.NewAccountService(accountRepo, transactionRepo)
	analyticsService := services.NewAnalyticsService(creditRepo, accountRepo, transactionRepo, paymentScheduleRepo)
	cardService := services.NewCardService(cardRepo)
	creditService := services.NewCreditService(paymentScheduleRepo, accountRepo, creditRepo)
	transactionService := services.NewTransactionService(db.DB, accountRepo, userRepo)

	userHandler := handlers.NewUserHandler(authService, utils.Logger)
	cardHandler := handlers.NewCardHandler(cardService, utils.Logger)
	accountHandler := handlers.NewAccountHandler(accountService, utils.Logger)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, utils.Logger)
	cbrHandler := handlers.NewCbrHandler(creditService, utils.Logger)
	creditHandler := handlers.NewCreditHandler(creditService, utils.Logger)
	transactionHandler := handlers.NewTransactionHandler(transactionService, utils.Logger)

	// Public routes
	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/cbr/keyrate", cbrHandler.GetKeyRate).Methods("GET")

	// Protected routes
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/accounts", accountHandler.GetAccounts).Methods("GET")
	protected.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	protected.HandleFunc("/accounts/{id}", accountHandler.UpdateAccount).Methods("PUT")
	protected.HandleFunc("/cards", cardHandler.GetCards).Methods("GET")
	protected.HandleFunc("/cards", cardHandler.IssueCard).Methods("POST")
	protected.HandleFunc("/transfer", transactionHandler.TransferFunds).Methods("POST")
	protected.HandleFunc("/analytics", analyticsHandler.GetAnalytics).Methods("GET")
	protected.HandleFunc("/credits", creditHandler.ApplyForCredit).Methods("POST")

	finalRouter := middleware.LoggingMiddleware(r)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", finalRouter))
}
