package deleteuser

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodDelete {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
				"success": false,
				"message": "Method not allowed",
			})
			return
		}

		id := r.PathValue("id")
		if id == "" {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "User ID is required",
			})
			return
		}

		// Check existence
		var exists string
		err := db.QueryRow("SELECT id FROM users WHERE id = $1", id).Scan(&exists)
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

		// Soft delete
		_, err = db.Exec("UPDATE users SET status = 0 WHERE id = $1", id)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": "Database error",
			})
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "User deactivated successfully",
		})
	}
}
