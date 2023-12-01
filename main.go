package main

import (
	"dapper/handlers"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	// Hardcoded data source name, normally would be in config file or secret
	dsn := "user=labs password=dapper host=localhost dbname=dapper sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	router := mux.NewRouter()

	router.HandleFunc("/login", handlers.LoginUser(db)).Methods("POST")

	router.HandleFunc("/users", handlers.GetUsers(db)).Methods("GET")
	router.HandleFunc("/users", handlers.UpdateUser(db)).Methods("PUT")

	router.HandleFunc("/signup", handlers.CreateUser(db)).Methods("POST")

	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		log.Fatalf("Error starting server: %v", err.Error())
	}
}
