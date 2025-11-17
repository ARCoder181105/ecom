package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func GetClaims(r *http.Request) (*Claims, error) {
	claims, ok := r.Context().Value(ClaimsContextKey).(*Claims)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	return claims, nil
}

func ParsePrice(priceStr string) float64 {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0.0
	}
	return price
}
