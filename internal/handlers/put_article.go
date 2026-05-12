package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/DilbarAkkaya/go-habr-microservices/internal/database"
	"github.com/DilbarAkkaya/go-habr-microservices/internal/dto/models"
	"github.com/DilbarAkkaya/go-habr-microservices/internal/middleware"

	"github.com/redis/go-redis/v9"
)

func UpdateArticleHandler(db *sql.DB, cache *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/articles/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		user, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		existing, err := database.GetArticle(db, id)
		if err != nil {
			http.Error(w, "article not found", http.StatusNotFound)
			return
		}

		if existing.AuthorID != user.ID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		var req models.ArticleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		article := &models.Article{
			ID:       id,
			AuthorID: user.ID,
			Title:    req.Title,
			Body:     req.Body,
		}

		if err := database.UpdateArticle(db, article); err != nil {
			http.Error(w, "could not update article", http.StatusInternalServerError)
			return
		}

		cache.Del(r.Context(), "articles:page:1")
		json.NewEncoder(w).Encode(article)
	}
}
