package database

import (
	"database/sql"
	"habr-app/internal/models"
)

func CreateUser(db *sql.DB, user *models.User) error {
	query := `INSERT INTO users (email, password_hash) VALUES ($1,$2) RETURNING id, created_at`
	err := db.QueryRow(query, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}