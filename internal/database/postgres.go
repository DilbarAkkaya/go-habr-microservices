package database

import (
	"database/sql"
	"fmt"
	"habr-app/internal/config"
	"os"

	_ "github.com/lib/pq"
)

func Connect()(*sql.DB, error){
	connStr:= config.GetDBConnStr()
	db, err:= sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf ("DB opening error: %w", err)
	}
	if err := db.Ping(); err!= nil {
		return nil, fmt.Errorf("DB is't answering to ping: %w", err)
	}
	return db, nil
}

func InitSchema(db *sql.DB) error{
	migrationQuery, err := os.ReadFile("migrations/000001_create_users_table.up.sql")
if err !=nil {
	return fmt.Errorf("Cannot read file migration %w", err)
}
_, err= db.Exec(string(migrationQuery))
if err !=nil {
	return fmt.Errorf("Cannot exec migration %w", err)
}
return nil
}