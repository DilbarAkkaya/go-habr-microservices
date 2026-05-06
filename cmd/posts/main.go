package main

import (
	"fmt"
	"habr-app/internal/config"
	"habr-app/internal/database"
	"habr-app/internal/handlers"
	"habr-app/internal/middleware"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC RECOVERED: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

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
	fmt.Println("Posts servcie started on :8082")
	http.Handle("/articles", middleware.JWTAuth(handlers.ArticlesHandler(db, redisClient)))
	http.Handle("/articles/", middleware.JWTAuth(handlers.ArticleHandler(db, redisClient)))
	port := config.GetServerPort("8082")
	log.Fatal(http.ListenAndServe(":"+port, recoverMiddleware(http.DefaultServeMux)))
}
