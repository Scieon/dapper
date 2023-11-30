package handlers

import (
	"dapper/models"
	"dapper/util"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

func GetUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}
