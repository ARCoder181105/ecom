package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
	}
}

func RespondWithError(w http.ResponseWriter, status int, err error) {

	RespondWithJSON(w, status, map[string]string{"error": err.Error()})
}

func ParseJson(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}
