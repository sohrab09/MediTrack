package deleteuser

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

func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Method check
		if r.Method != http.MethodDelete {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
				"success": false,
				"message": "Method not allowed",
			})
			return
		}

		// Get ID from query param
		id := r.URL.Query().Get("id")
		if id == "" {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "User ID is required",
			})
			return
		}

		// DB Query - check if user exists
		var user models.User
		err := db.QueryRow("SELECT id FROM users WHERE id = $1", id).Scan(&user.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				respondJSON(w, http.StatusNotFound, map[string]interface{}{
					"success": false,
					"message": "User not found",
				})
				return
			}
			respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": "Database error",
			})
			return
		}

		// DB Query - soft delete: set status = 0

		_, err = db.Exec("UPDATE users SET status = 0 WHERE id = $1", id)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": "Database error",
			})
			return
		}

		// Send success response
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "User deactivated successfully",
		})
	}
}
