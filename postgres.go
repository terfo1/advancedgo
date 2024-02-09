package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func authenticateUser(username, password string) bool {
	const connStr = "user=postgres password=terfo2005 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var storedPassword string
	err = db.QueryRow(`SELECT "password" FROM clients WHERE "user" = $1`, username).Scan(&storedPassword)
	if err == sql.ErrNoRows {
		// User not found in the database
		return false
	} else if err != nil {
		log.Fatal(err)
	}

	return storedPassword == password
}

func getUsersFromDB() ([]User, error) {
	const connStr = "user=postgres password=terfo2005 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("Error opening database connection: %s", err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT "user" FROM clients`)
	if err != nil {
		return nil, fmt.Errorf("Error querying database: %s", err)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		err := rows.Scan(&u.Username)
		if err != nil {
			return nil, fmt.Errorf("Error scanning row: %s", err)
		}
		users = append(users, u)
	}

	return users, nil
}
func userExists(username string) bool {
	const connStr = "user=postgres password=terfo2005 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database connection: %s\n", err)
		return false
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM clients WHERE "user" = $1)`, username).Scan(&exists)
	if err != nil {
		log.Printf("Error querying user existence: %s\n", err)
		return false
	}

	return exists
}
func insertUser(user User) error {
	const connStr = "user=postgres password=terfo2005 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("Error opening database connection: %s", err)
	}
	defer db.Close()
	_, err = db.Exec(`INSERT INTO clients ("user", "password") VALUES ($1, $2)`, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("Error inserting user into database: %s", err)
	}

	return nil
}
func insertdata(data Profile) error {
	fmt.Println("data is here", data)
	const connStr = "user=postgres password=terfo2005 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("Error opening database connection: %s", err)
	}
	defer db.Close()
	parsedTime, err := time.Parse("2006-01-02", data.Dateofbirthyear)
	if err != nil {
		return fmt.Errorf("Invalid date format", err)
	}
	_, err = db.Exec(`UPDATE clients SET imya = $1, surname = $2, dateofbirth = $3 WHERE user=$4`, data.Name, data.Surname, parsedTime, data.Username)
	if err != nil {
		return fmt.Errorf("Error inserting user into database: %s", err)
	}

	return nil
}
