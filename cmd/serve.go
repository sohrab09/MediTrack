package cmd

import (
	"fmt"
	"log"
	global_router "meditrack-backend/golbel_router"
	"meditrack-backend/internal/config"
	"meditrack-backend/internal/database"
	"meditrack-backend/internal/handlers/getusers"
	"meditrack-backend/internal/handlers/login"
	"meditrack-backend/internal/handlers/register"
	"net/http"
	"time"
)

func Serve() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to DB
	db := database.ConnectPostgres(cfg)
	defer db.Close()

	// Create router
	mux := http.NewServeMux()
	globalHandler := global_router.GlobalRouter(mux)

	// Routes
	mux.HandleFunc("/api/v1/auth/login", login.LoginHandler(db))
	mux.HandleFunc("/api/v1/auth/register", register.RegisterHandler(db))
	mux.HandleFunc("/api/v1/users", getusers.GetUsers(db))

	// Start server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      globalHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 MediTrack server running on port %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
