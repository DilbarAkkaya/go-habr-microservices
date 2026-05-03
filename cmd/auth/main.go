package main

import (
	"fmt"
	"habr-app/internal/database"
	"habr-app/internal/handlers"
	"log"
	"net/http"
)

func main(){
	fmt.Println("Habr-app is starting...")
	db, err:= database.Connect()
	
	fmt.Println("Habr-app is connected to db...")
	err = database.InitSchema(db)
	if err != nil {
		log.Fatal("Error", err)
	}
		defer db.Close()
	http.HandleFunc("/register", handlers.RegisterHandler(db))
fmt.Println("Server started at http://localhost:8080")
log.Fatal(http.ListenAndServe(":8080", nil))
}