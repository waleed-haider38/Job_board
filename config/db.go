package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:waleedhaider@localhost:5432/jobboard?sslmode=disable")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database not reachable:", err)
	}

	log.Println("Database connected successfully!")
	return db
}
