package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DilbarAkkaya/go-habr-microservices/internal/config"
	"github.com/DilbarAkkaya/go-habr-microservices/internal/database"
	"github.com/DilbarAkkaya/go-habr-microservices/internal/handlers"
	"github.com/DilbarAkkaya/go-habr-microservices/internal/middleware"
)

func main() {
	fmt.Println("Auth service starting...")

	db, err := database.Connect()
	if err != nil {
		log.Fatal("db connect:", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlers.RegisterHandler(db))
	mux.HandleFunc("/login", handlers.LoginHandler(db))
	mux.HandleFunc("/verify", handlers.VerifyEmailHandler(db))

	port := config.GetServerPort("8081")
	log.Println("Server listening on :" + port)

	log.Fatal(http.ListenAndServe(":"+port, middleware.RecoverMiddleware(mux)))
}
