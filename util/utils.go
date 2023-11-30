package util

import (
	"dapper/models"
	"encoding/json"
	"net/http"
)

func SendErrorResponse(w http.ResponseWriter, r *http.Request, code int, message string, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errorResponse := models.ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}

	json.NewEncoder(w).Encode(errorResponse)
}
