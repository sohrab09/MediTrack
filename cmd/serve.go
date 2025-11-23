package cmd

import (
	"log"
	global_router "meditrack-backend/golbel_router"
	"meditrack-backend/internal/config"
	"meditrack-backend/internal/database"
	deleteuser "meditrack-backend/internal/handlers/delete-user"
	getmedicinecategories "meditrack-backend/internal/handlers/get_medicine_categories"
	"meditrack-backend/internal/handlers/getuser"
	"meditrack-backend/internal/handlers/getusers"
	"meditrack-backend/internal/handlers/login"
	medicinecategories "meditrack-backend/internal/handlers/medicine_categories"
	"meditrack-backend/internal/handlers/register"
	"net/http"
	"time"
)

// func Serve() {
// 	// Load config
// 	cfg := config.LoadConfig()

// 	// Connect to DB
// 	db := database.ConnectPostgres(cfg)
// 	defer db.Close()

// 	// Create router
// 	mux := http.NewServeMux()
// 	globalHandler := global_router.GlobalRouter(mux)

// 	// Routes
// 	mux.HandleFunc("/api/v1/auth/login", login.LoginHandler(db))
// 	mux.HandleFunc("/api/v1/auth/register", register.RegisterHandler(db))
// 	mux.HandleFunc("/api/v1/users", getusers.GetUsers(db))
// 	mux.HandleFunc("/api/v1/user", getuser.GetUser(db))
// 	mux.HandleFunc("/api/v1/user/delete", deleteuser.DeleteUser(db))

// 	// Routes for Medicine Categories
// 	mux.HandleFunc("POST /api/v1/medicine-categories", medicinecategories.CreateMedicineCategories(db))
// 	mux.HandleFunc("GET /api/v1/medicine-categories", getmedicinecategories.GetMedicineCategories(db))

// 	// Start server
// 	srv := &http.Server{
// 		Addr:         fmt.Sprintf(":%s", cfg.Port),
// 		Handler:      globalHandler,
// 		ReadTimeout:  10 * time.Second,
// 		WriteTimeout: 15 * time.Second,
// 		IdleTimeout:  60 * time.Second,
// 	}

// 	log.Printf("🚀 MediTrack server running on port %s", cfg.Port)
// 	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 		log.Fatalf("Server error: %v", err)
// 	}
// }

func Serve() {
	cfg := config.LoadConfig()
	db := database.ConnectPostgres(cfg)
	defer db.Close()

	mux := http.NewServeMux()
	globalHandler := global_router.GlobalRouter(mux)

	// Auth
	mux.HandleFunc("POST /api/v1/auth/login", login.LoginHandler(db))
	mux.HandleFunc("POST /api/v1/auth/register", register.RegisterHandler(db))

	// Users
	mux.HandleFunc("GET /api/v1/users", getusers.GetUsers(db))
	mux.HandleFunc("GET /api/v1/users/{id}", getuser.GetUser(db))
	mux.HandleFunc("DELETE /api/v1/users/{id}", deleteuser.DeleteUser(db))

	// Medicine Categories
	mux.HandleFunc("POST /api/v1/medicine-categories", medicinecategories.CreateMedicineCategories(db))
	mux.HandleFunc("GET /api/v1/medicine-categories", getmedicinecategories.GetMedicineCategories(db))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
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
