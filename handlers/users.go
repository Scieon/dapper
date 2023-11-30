package handlers

import (
	"dapper/models"
	"dapper/util"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var secretKey = []byte("dapper")

func GetUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the JWT from the header
		tokenString := r.Header.Get("X-Authentication-Token")

		email := validateJWT(tokenString)

		if email == "" {
			util.SendErrorResponse(w, r, http.StatusUnauthorized, "Invalid authentication token", "")
			return
		}

		var users []models.User

		result := db.Find(&users)
		if result.Error != nil {
			panic("failed to query users")
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(users)

		if err != nil {
			util.SendErrorResponse(w, r, http.StatusInternalServerError, "Invalid request", err.Error())
			return
		}
	}
}

func CreateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			util.SendErrorResponse(w, r, http.StatusBadRequest, "Invalid request", err.Error())
			return
		}

		result := db.Create(&user)
		if result.Error != nil {
			util.SendErrorResponse(w, r, http.StatusBadRequest, "Invalid request", result.Error.Error())
			return
		}

		jwtToken, err := createJWT(user.Email)
		if err != nil {
			util.SendErrorResponse(w, r, http.StatusInternalServerError, "Could not create JWT", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(jwtToken)
	}
}

func LoginUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var login models.LoginBody

		err := json.NewDecoder(r.Body).Decode(&login)
		if err != nil {
			util.SendErrorResponse(w, r, http.StatusBadRequest, "Invalid request", err.Error())
			return
		}

		var user models.User
		_ = db.Where("email = ? AND password = ?", login.Email, login.Password).First(&user)

		// Check if user exists
		if user.Email == "" {
			util.SendErrorResponse(w, r, http.StatusUnauthorized, "Invalid credentials", "")
			return
		}

		// Ignore err - for cleaner code
		jwtToken, _ := createJWT(login.Email)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwtToken)
	}
}

func UpdateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the JWT from the header
		tokenString := r.Header.Get("X-Authentication-Token")

		// Get email from token
		email := validateJWT(tokenString)

		if email == "" {
			util.SendErrorResponse(w, r, http.StatusUnauthorized, "Invalid authentication token", "")
			return
		}

		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			util.SendErrorResponse(w, r, http.StatusBadRequest, "Invalid request", err.Error())
			return
		}

		fmt.Printf("===> %s", email)

		// build query string slowly - fn, ln, fn ln
		db = db.Debug()
		_ = db.Model(user).Where("email = ?", email).
			Updates(map[string]interface{}{"firstname": user.FirstName, "lastname": user.LastName})

		w.Header().Set("Content-Type", "application/json")
	}
}

func validateJWT(tokenString string) string {
	// Parse and validate the JWT
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ignore signing method and just return secret
		return secretKey, nil
	})

	if err != nil || !jwtToken.Valid {
		return ""
	}

	// Token is valid; you can access its claims
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Failed to extract claims from JWT")
		return ""
	}

	// Extract user-related information from the claims
	email, ok := claims["user_email"].(string)

	return email
}

func createJWT(email string) (string, error) {
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
