package cmd

import (
	"fmt"
	"log"
	"meditrack-backend/internal/config"
	"meditrack-backend/internal/database"
	"net/http"
	"time"
)

func Serve() {
	// Load config
	cfg := config.LoadConfig()

	// Connect DB
	db := database.ConnectPostgres(cfg)
	// keep db globally available to handlers package (we'll set it inside database package)
	_ = db
	defer db.Close()

	// Routes
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 MediTrack server running on port %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
