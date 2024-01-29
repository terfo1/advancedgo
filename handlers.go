package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Profile struct {
	Name             string `json:"name"`
	Surname          string `json:"surname"`
	Dateofbirthyear  string `json:"dateofbirthyear"`
	Dateofbirthmonth string `json:"dateofbirthmonth"`
	Dateofbirthday   string `json:"dateofbirthday"`
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postreg(w, r)
		return
	}

	if r.Method == http.MethodGet && r.URL.Path == "/users" {
		getUsersHandler(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
func edit_Profile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ts, err := template.ParseFiles("edit_profile.html")
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = ts.Execute(w, nil)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	case http.MethodPost:
		edit(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func edit(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %s", err), http.StatusInternalServerError)
		return
	}
	var data Profile
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON data: %s", err), http.StatusBadRequest)
		return
	}
	err = insertdata(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting user into database: %s", err), http.StatusInternalServerError)
		return
	}
}
func insertdata(data Profile) error {
	const connStr = "user=postgres password=terfo2005 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("Error opening database connection: %s", err)
	}
	defer db.Close()
	dob := fmt.Sprintf("%s-%s-%s", data.Dateofbirthyear, data.Dateofbirthmonth, data.Dateofbirthday)
	_, err = db.Exec(`INSERT INTO clients ("imya", "surname","dateofbirth") VALUES ($1, $2,$3)`, data.Name, data.Surname, dob)
	if err != nil {
		return fmt.Errorf("Error inserting user into database: %s", err)
	}

	return nil
}
func postreg(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %s", err), http.StatusInternalServerError)
		return
	}

	var incomingUser User
	err = json.Unmarshal(body, &incomingUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON data: %s", err), http.StatusBadRequest)
		return
	}

	err = insertUser(incomingUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting user into database: %s", err), http.StatusInternalServerError)
		return
	}
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
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := getUsersFromDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving users from database: %s", err), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("users").Parse(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>User List</title>
		</head>
		<body>
			<h2>User List</h2>
			<ul>
				{{range .}}
					<li>{{.Username}}</li>
				{{end}}
			</ul>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, "Error rendering HTML", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, users)
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ts, err := template.ParseFiles("Login.html")
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = ts.Execute(w, nil)
	case http.MethodPost:
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		fmt.Println("Received login request for username:", user.Username)

		if authenticateUser(user.Username, user.Password) {
			fmt.Println("Login successful for username:", user.Username)
			http.Redirect(w, r, "/edit_profile?username="+url.QueryEscape(user.Username), http.StatusSeeOther)
			return
		} else {
			fmt.Println("Login failed for username:", user.Username)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
		}

	}
}

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

	// Compare the stored password with the provided password
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
