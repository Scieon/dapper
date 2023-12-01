package util

import (
	"dapper/models"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

var secretKey = []byte("dapper")

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

func ValidateJWT(tokenString string) *string {
	// Parse and validate the JWT
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ignore signing method and just return secret
		return secretKey, nil
	})

	if err != nil || !jwtToken.Valid {
		return nil
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Failed to extract claims from JWT")
		return nil
	}

	// Extract user-related information from the claims
	email, ok := claims["user_email"].(string)

	return &email
}

func CreateJWT(email string) (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)

	claims := jwt.MapClaims{
		"user_email": email,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	jwtToken.Claims = claims

	// Sign token
	tokenString, err := jwtToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
