package handlers

import (
	"database/sql"
	"encoding/json"
	"habr-app/internal/database"
	"habr-app/internal/middleware"
	"habr-app/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type articleRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func ArticlesHandler(db *sql.DB, cache *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("ArticlesHandler called for %s", r.Method)
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		switch r.Method {
		case http.MethodGet:
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
		case http.MethodPost:
			var req articleRequest
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
}

func ArticleHandler(db *sql.DB, cache *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/articles/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			art, err := database.GetArticle(db, id)
			if err != nil {
				http.Error(w, "article not found", http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(art)
		case http.MethodPut:
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

			var req articleRequest
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

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
