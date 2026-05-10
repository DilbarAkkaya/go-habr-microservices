package database

import (
	"database/sql"
	"habr-app/internal/models"
	"log"
	"time"
)

func CreateArticle(db *sql.DB, article *models.Article) error {
	query := `
	    INSERT INTO articles (author_id, title, body, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	log.Printf("Executing query: %s\n", query)
	log.Printf("Params: author_id=%d, title=%s, body=%s\n", article.AuthorID, article.Title, article.Body)
	now := time.Now()

	return db.QueryRow(query, article.AuthorID, article.Title, article.Body, now, now).Scan(&article.ID, &article.CreatedAt, &article.UpdatedAt)

}

func UpdateArticle(db *sql.DB, article *models.Article) error {
	query := `
UPDATE articles
SET title = $1, body = $2, updated_at = $3
WHERE id = $4
RETURNING author_id, created_at, updated_at
`
	now := time.Now()
	return db.QueryRow(query, article.Title, article.Body, now, article.ID).Scan(&article.AuthorID, &article.CreatedAt, &article.UpdatedAt)
}
func GetArticle(db *sql.DB, id int) (*models.Article, error) {
	query := `
	    SELECT id, author_id, title, body, created_at, updated_at
		FROM articles
		WHERE id = $1
	`
	art := &models.Article{}
	err := db.QueryRow(query, id).Scan(&art.ID, &art.AuthorID, &art.Title, &art.Body, &art.CreatedAt, &art.UpdatedAt)
	return art, err
}

func ListArticles(db *sql.DB, limit, offset int) ([]models.Article, error) {
	query := `
       SELECT id, author_id, title, body, created_at, updated_at
	   FROM articles
	   ORDER BY created_at DESC
	   LIMIT $1 OFFSET $2  
   `
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var art models.Article
		if err := rows.Scan(&art.ID, &art.AuthorID, &art.Title, &art.Body, &art.CreatedAt, &art.UpdatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, art)
	}
	return articles, nil

}
