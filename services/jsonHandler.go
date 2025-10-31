package respond

import (
	"encoding/json"
	"net/http"
)


func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
	}
}

// RespondWithError provides a unified error response structure.
func RespondWithError(w http.ResponseWriter, status int, err interface{}) {
	response := map[string]interface{}{
		"status":  "error",
		"code":    status,
		"message": err,
	}
	RespondWithJSON(w, status, response)
}
