package database

import (
	"database/sql"
	"log"
	"meditrack-backend/internal/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectPostgres(cfg *config.Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}

	// simple ping to verify
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
	}

	// set global
	DB = db
	log.Println("✅ Connected to PostgreSQL")
	return db
}
