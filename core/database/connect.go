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

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type ContextKey string

const UserContextKey ContextKey = "user"

func QueryUser(email string) error {
	rows, err := DB.Query("SELECT id, email FROM users WHERE email = $1", email)
	if err != nil {
		return fmt.Errorf("Err executing select query: %v\n", err)
	}

	for rows.Next() {
		var id int
		var email string
		err := rows.Scan(&id, &email)
		if err != nil {
			return fmt.Errorf("Err scanning next row: %v\n", err)
		}
		fmt.Printf("id: %d, name: %s\n", id, email)
	}

	// Check for errors from iterating over rows
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("Err iterating over rows: %v\n", err)
	}

	return nil
}

func QueryToken(token string) (*User, error) {
	var user User

	row := DB.QueryRow("SELECT id, email FROM users WHERE token = $1", token)
	err := row.Scan(&user.Id, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return &user, fmt.Errorf("User not found")
		} else {
			log.Fatalf("QueryToken: %v", err)
		}
	}

	return &user, nil
}
