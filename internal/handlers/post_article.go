package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/DilbarAkkaya/go-habr-microservices/internal/database"
	"github.com/DilbarAkkaya/go-habr-microservices/internal/dto/models"
	"github.com/DilbarAkkaya/go-habr-microservices/internal/middleware"

	"github.com/redis/go-redis/v9"
)

func CreateArticleHandler(db *sql.DB, cache *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("CreateArticleHandler called")
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.ArticleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		user, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		log.Printf("Creating article userID = %d | email = %s\n", user.ID, user.Email)
		article := &models.Article{
			AuthorID: user.ID,
			Title:    req.Title,
			Body:     req.Body,
		}
		if err := database.CreateArticle(db, article); err != nil {
			http.Error(w, "could not create article", http.StatusInternalServerError)
			return
		}

		cache.Del(r.Context(), "articles:page:1")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(article)
	}
}
