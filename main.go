package main

import (
	"dapper/handlers"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	// Hardcoded data source name, normally would be in config file or secret
	dsn := "host=dapper-pg user=labs password=dapper dbname=dapper port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Hello World")
	}).Methods("GET")

	router.HandleFunc("/login", handlers.LoginUser(db)).Methods("POST")

	router.HandleFunc("/users", handlers.GetUsers(db)).Methods("GET")
	router.HandleFunc("/users", handlers.UpdateUser(db)).Methods("PUT")

	router.HandleFunc("/signup", handlers.CreateUser(db)).Methods("POST")

	fmt.Println("Server is started on localhost:8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Error starting server: %v", err.Error())
	}
}
