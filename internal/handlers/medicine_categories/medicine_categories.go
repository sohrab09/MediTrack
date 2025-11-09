package medicinecategories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"meditrack-backend/internal/models"
	"net/http"
)

// Helper to send JSON response
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func CreateMedicineCategories(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Method check
		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
				"success": false,
				"message": "Method not allowed",
			})
			return
		}

		// Decode request
		var req models.MedicineCategories
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Invalid JSON format",
			})
			return
		}

		// Taking input
		var data models.MedicineCategories

		data = req

		fmt.Println(data)

		// validate input
		if data.Name == "" || data.Description == "" || data.Status == 0 {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "All fields are required",
			})
			return
		}

		// Duplication check
		var category models.MedicineCategories
		err := db.QueryRow("SELECT id FROM categories WHERE name = $1", data.Name).Scan(&category.ID)
		if err == nil {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Category already exists",
			})
			return
		}

		// DB Query
		_, err = db.Exec("INSERT INTO categories (name, description, status) VALUES ($1, $2, $3)", data.Name, data.Description, data.Status)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": "Database error",
			})
			return
		}

		// Success response
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Category created successfully",
		})
	}
}
