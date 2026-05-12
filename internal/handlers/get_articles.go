package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/DilbarAkkaya/go-habr-microservices/internal/database"
	"github.com/DilbarAkkaya/go-habr-microservices/internal/dto/models"

	"github.com/redis/go-redis/v9"
)

func ListArticlesHandler(db *sql.DB, cache *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("ListArticlesHandler called")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		ctx := r.Context()
		cacheKey := "articles:page:1"
		if page == 1 {
			if cached, err := cache.Get(ctx, cacheKey).Result(); err == nil {
				var list []models.Article
				if err := json.Unmarshal([]byte(cached), &list); err == nil {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(list)
					return
				}
			}
		}
		articles, err := database.ListArticles(db, 10, (page-1)*10)
		if err != nil {
			http.Error(w, "could not list articles", http.StatusInternalServerError)
			return
		}
		if page == 1 {
			data, _ := json.Marshal(articles)
			cache.Set(ctx, cacheKey, data, 30*time.Second)
		}
		json.NewEncoder(w).Encode(articles)
	}
}
