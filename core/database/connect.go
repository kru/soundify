package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	// Open a connection to the database
	var err error
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if len(dbHost) == 0 || len(dbPort) == 0 || len(dbUser) == 0 || len(dbPassword) == 0 || len(dbName) == 0 {
		log.Fatalf("Database credential is not set")
	}

	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error ping the database %v", err)
	}
	fmt.Println("Successfully connected to the database!")
}
