package db

import (
	"database/sql"
	"fmt"
	"os"
	
	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DATABASE")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, database))
	if err != nil {
		return nil, err
	}

	return db, nil
}
