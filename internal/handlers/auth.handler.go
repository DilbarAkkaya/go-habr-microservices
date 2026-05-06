package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"habr-app/internal/broker"
	"habr-app/internal/database"
	"habr-app/internal/models"
	"habr-app/internal/utils"
	"log"
	"net/http"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode((&req)); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Password handler error", http.StatusInternalServerError)
			return
		}
		user := &models.User{
			Email:             req.Email,
			PasswordHash:      hashedPassword,
			VerificationToken: utils.RandomString(32),
		}
		err = database.CreateUser(db, user)
		if err != nil {
			if errors.Is(err, database.ErrUserAlreadyExists) {
				http.Error(w, "user already exists", http.StatusConflict)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if err := broker.SendRegistrationEvent(user.Email, user.VerificationToken); err != nil {
			log.Println("kafka error:", err)

		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "user created, please verify email",
		})
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return

		}
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		user, err := database.GetUserByEmail(db, req.Email)
		if err != nil {
			http.Error(w, "invalid credentals", http.StatusUnauthorized)
			return
		}
		log.Printf("DEBUG: DB_HASH: [%s] | REQ_PASS: [%s]", user.PasswordHash, req.Password)
		if !utils.CheckHashPassword(req.Password, user.PasswordHash) {
			http.Error(w, "invalid credentals", http.StatusUnauthorized)
			return
		}
		if !user.EmailConfirmed {
			http.Error(w, "email not confirmed", http.StatusForbidden)
			return
		}

		token, err := utils.GenerateJWT(user.ID, user.Email)
		if err != nil {
			http.Error(w, "could not generate token", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{
			"access_token": token,
		})
	}
}

func VerifyEmailHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "missing token", http.StatusBadRequest)
			return
		}

		if err := database.ConfirmUserEmail(db, token); err != nil {
			http.Error(w, "invalid token", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{
			"message": "email confirmed",
		})
	}
}
