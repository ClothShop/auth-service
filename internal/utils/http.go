package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error, message string, statusCode int) {
	log.Printf("❌ %s: %v", message, err)
	http.Error(w, message, statusCode)
}

func DecodeJSONRequest[T any](w http.ResponseWriter, r *http.Request, dst *T) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		HandleError(w, err, "❌ Invalid JSON format", http.StatusBadRequest)
		return false
	}
	return true
}
