package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/DilbarAkkaya/go-habr-microservices/internal/database"

	"github.com/redis/go-redis/v9"
)

func GetArticleHandler(db *sql.DB, cache *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/articles/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		art, err := database.GetArticle(db, id)
		if err != nil {
			http.Error(w, "article not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(art)
	}
}
