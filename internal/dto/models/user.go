package models

import "time"

type User struct {
	ID                int       `json:"id"`
	Email             string    `json:"email"`
	PasswordHash      string    `json:"-"`
	EmailConfirmed    bool      `json:"email_confirmed"`
	VerificationToken string    `json:"-"`
	CreatedAt         time.Time `json:"created_at"`
}
