package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DilbarAkkaya/go-habr-microservices/internal/config"
	"github.com/DilbarAkkaya/go-habr-microservices/internal/database"
	"github.com/DilbarAkkaya/go-habr-microservices/internal/handlers"
	"github.com/DilbarAkkaya/go-habr-microservices/internal/middleware"

	"github.com/redis/go-redis/v9"
)

func main() {
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET not set")
	}
	db, err := database.Connect()
	if err != nil {
		log.Fatal("db connect:", err)
	}
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.GetRedisAddr(),
	})
	defer redisClient.Close()
	port := config.GetServerPort("8082")
	fmt.Println("Posts servcie started on :", port)
	http.Handle("/articles", middleware.JWTAuth(handlers.ListArticlesHandler(db, redisClient)))
	http.Handle("/articles/", middleware.JWTAuth(handlers.GetArticleHandler(db, redisClient)))

	log.Fatal(http.ListenAndServe(":"+port, middleware.RecoverMiddleware(http.DefaultServeMux)))
}
