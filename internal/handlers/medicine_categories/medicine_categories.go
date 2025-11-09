package medicinecategories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"meditrack-backend/internal/models"
	"net/http"
	"strings"
	"time"
)

// Response is the common API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Helper to send consistent JSON responses
func respondJSON(w http.ResponseWriter, status int, res Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

// Validate category input fields
func validateCategoryInput(data *models.MedicineCategories) error {
	if strings.TrimSpace(data.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(data.Description) == "" {
		return errors.New("description is required")
	}
	if data.Status != 0 && data.Status != 1 {
		return errors.New("status must be 0 (inactive) or 1 (active)")
	}
	return nil
}

func CreateMedicineCategories(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Method check
		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, Response{
				Success: false,
				Message: "Method not allowed",
			})
			return
		}

		ctx := r.Context()
		var req struct {
			Data models.MedicineCategories `json:"data"`
		}

		// Decode JSON request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, Response{
				Success: false,
				Message: "Invalid JSON format",
			})
			return
		}
		data := req.Data

		// Validate input
		if err := validateCategoryInput(&data); err != nil {
			respondJSON(w, http.StatusBadRequest, Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		// Check for duplicates
		var existingID int
		err := db.QueryRowContext(ctx, "SELECT id FROM categories WHERE name = $1", data.Name).Scan(&existingID)
		if err != nil && err != sql.ErrNoRows {
			log.Println("Database error on duplicate check:", err)
			respondJSON(w, http.StatusInternalServerError, Response{
				Success: false,
				Message: "Database query error",
			})
			return
		}
		if err == nil {
			respondJSON(w, http.StatusBadRequest, Response{
				Success: false,
				Message: "Category already exists",
			})
			return
		}

		// Insert into DB
		query := `
			INSERT INTO categories (name, description, status, created_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at
		`

		var createdID int
		var createdAt time.Time
		err = db.QueryRowContext(ctx, query, data.Name, data.Description, data.Status, time.Now()).Scan(&createdID, &createdAt)
		if err != nil {
			log.Println("Database insert error:", err)
			respondJSON(w, http.StatusInternalServerError, Response{
				Success: false,
				Message: "Failed to create category",
			})
			return
		}

		data.ID = createdID
		data.CreatedAt = createdAt.Format(time.RFC3339)

		log.Printf("Category created successfully: %+v", data)

		// Respond with created data
		respondJSON(w, http.StatusCreated, Response{
			Success: true,
			Message: "Category created successfully",
			Data:    data,
		})
	}
}
