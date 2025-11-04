package getusers

import (
	"database/sql"
	"encoding/json"
	"meditrack-backend/internal/models"
	"net/http"
)

// Helper to send JSON response
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Method check
		if r.Method != http.MethodGet {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
				"success": false,
				"message": "Method not allowed",
			})
			return
		}

		// DB Query
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": "Database error",
			})
			return
		}
		defer rows.Close()

		// Data
		var users []models.User
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Phone, &user.Email, &user.Password, &user.CreatedAt); err != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
					"success": false,
					"message": "Database error",
				})
				return
			}
			users = append(users, user)
		}
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"data":    users,
			"count":   len(users),
		})
	}
}
