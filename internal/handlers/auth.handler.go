package handlers

import (
	"database/sql"
	"encoding/json"
	"habr-app/internal/database"
	"habr-app/internal/models"
	"habr-app/internal/utils"
	"net/http"
)

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err:= json.NewDecoder(r.Body).Decode((&user)); err!=nil {
http.Error(w, "Invalid request", http.StatusBadRequest)
return
		}
		hashedPassword, err:= utils.HashPassword(user.PasswordHash)
		if err != nil {
			http.Error(w, "Password handler error", http.StatusInternalServerError)
			return
		}
		user.PasswordHash = hashedPassword
		err = database.CreateUser(db, &user)
		if err != nil {
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}