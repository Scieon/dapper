package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var secretKey = []byte("dapper")

func ValidateJWT(tokenString string) *string {
	// Parse and validate the JWT
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ignore signing method and just return secret
		return secretKey, nil
	})

	if err != nil || !jwtToken.Valid {
		return nil
	}

	// Token is valid; you can access its claims
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
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (adjust as needed)
	}

	jwtToken.Claims = claims

	// Sign token
	tokenString, err := jwtToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
