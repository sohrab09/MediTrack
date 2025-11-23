package getmedicinecategories

import (
	"database/sql"
	"encoding/json"
	"meditrack-backend/internal/models"
	"net/http"
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

func GetMedicineCategories(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Method check
		if r.Method != http.MethodGet {
			respondJSON(w, http.StatusMethodNotAllowed, Response{
				Success: false,
				Message: "Method not allowed",
			})
			return
		}

		// Query specific columns, excluding password
		rows, err := db.Query("SELECT id, name, description, status, created_at FROM categories")
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, Response{
				Success: false,
				Message: "Database error",
			})
			return
		}
		defer rows.Close()

		var categories []models.MedicineCategories
		for rows.Next() {
			var category models.MedicineCategories
			if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.Status, &category.CreatedAt); err != nil {
				respondJSON(w, http.StatusInternalServerError, Response{
					Success: false,
					Message: "Database error",
				})
				return
			}
			categories = append(categories, category)
		}

		respondJSON(w, http.StatusOK, Response{
			Success: true,
			Message: "Categories fetched successfully",
			Data:    categories,
		})
	}
}
