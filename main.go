package main

import (
	"dapper/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	dsn := "user=labs password=dapper host=localhost dbname=dapper sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {

	}

	router := mux.NewRouter()

	router.HandleFunc("/login", handlers.LoginUser(db)).Methods("POST")

	router.HandleFunc("/users", handlers.GetUsers(db)).Methods("GET")

	router.HandleFunc("/signup", handlers.CreateUser(db)).Methods("POST")

	fmt.Println("server starting on 8080")
	http.ListenAndServe("localhost:8080", router)
}
