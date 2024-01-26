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
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("index.html")
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = ts.Execute(w, nil)
	})

	http.HandleFunc("/register", registrationHandler)
	http.HandleFunc("/users", getUsersHandler)

	port := 8080
	fmt.Printf("Server is running on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
