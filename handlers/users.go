package handlers

import (
	"dapper/models"
	"dapper/util"
	"encoding/json"
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

		if tokenString == "" {
			util.SendErrorResponse(w, r, http.StatusUnauthorized, "Missing authentication token", "")
			return
		}

		// Parse and validate the JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ignore signing method and just return secret
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			util.SendErrorResponse(w, r, http.StatusUnauthorized, "Invalid authentication token", err.Error())
			return
		}

		var users []models.User

		result := db.Find(&users)
		if result.Error != nil {
			panic("failed to query users")
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(users)

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

		token, err := createJWT(user.Email)
		if err != nil {
			util.SendErrorResponse(w, r, http.StatusInternalServerError, "Could not create JWT", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(token)
	}
}

func createJWT(email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := jwt.MapClaims{
		"user_email": email,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (adjust as needed)
	}

	token.Claims = claims

	// Sign token
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
