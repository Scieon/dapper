package handlers

import (
	"dapper/models"
	"dapper/token"
	"dapper/util"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

func GetUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the JWT from the header
		tokenString := r.Header.Get("X-Authentication-Token")

		email := token.ValidateJWT(tokenString)

		if email == nil {
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

		// Check request is valid
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

		// Check is jwt valid
		jwtToken, err := token.CreateJWT(user.Email)
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
		jwtToken, _ := token.CreateJWT(login.Email)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwtToken)
	}
}

func UpdateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the JWT from the header
		tokenString := r.Header.Get("X-Authentication-Token")

		// Get email from token
		email := token.ValidateJWT(tokenString)

		if email == nil {
			util.SendErrorResponse(w, r, http.StatusUnauthorized, "Invalid authentication token", "")
			return
		}

		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			util.SendErrorResponse(w, r, http.StatusBadRequest, "Invalid request", err.Error())
			return
		}

		// Create a map to hold the updates
		updates := map[string]interface{}{}

		// Check if the last name is present and not empty before adding it to updates
		// should also check to see if neither exists...
		if user.FirstName != "" {
			updates["firstname"] = user.FirstName
		}

		if user.LastName != "" {
			updates["lastname"] = user.LastName
		}

		db = db.Debug()
		_ = db.Model(user).Where("email = ?", email).Updates(updates)

		w.Header().Set("Content-Type", "application/json")
	}
}
