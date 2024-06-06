package database

import (
	"database/sql"
	"fmt"
	"log"
)

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
