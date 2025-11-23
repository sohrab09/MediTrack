package getuser

import (
	"database/sql"
	"encoding/json"
	"meditrack-backend/internal/models"
	"net/http"
)

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
				"success": false,
				"message": "Method not allowed",
			})
			return
		}

		// Extract `{id}` from URL path
		id := r.PathValue("id")
		if id == "" {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "User ID is required",
			})
			return
		}

		var user models.User
		err := db.QueryRow(
			"SELECT id, first_name, last_name, phone, email, created_at FROM users WHERE id=$1",
			id,
		).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Phone, &user.Email, &user.CreatedAt)

		if err == sql.ErrNoRows {
			respondJSON(w, http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": "User not found",
			})
			return
		}

		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": "Database error",
			})
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"data":    user,
		})
	}
}
