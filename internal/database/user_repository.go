package database

import (
	"database/sql"
	"errors"
	"habr-app/internal/models"
	"strings"
)

var ErrUserAlreadyExists = errors.New("user already exists")

func CreateUser(db *sql.DB, user *models.User) error {
	query := `
	    INSERT INTO users (email, password_hash, verification_token)
	    VALUES ($1,$2, $3)
		RETURNING id, created_at
	`
	err := db.QueryRow(query, user.Email, user.PasswordHash, user.VerificationToken).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			return ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	query := `
	SELECT id, email, password_hash, email_confirmed, verification_token, created_at
	FROM users
	WHERE email=$1
	`
	user := &models.User{}
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.EmailConfirmed, &user.VerificationToken, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func ConfirmUserEmail(db *sql.DB, token string) error {
	query := `
	    UPDATE users 
	    SET email_confirmed = true, verification_token = ''
		WHERE verification_token = $1
		`
	res, err := db.Exec(query, token)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Invalid token")
	}
	return nil
}
