package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	_ "github.com/sirupsen/logrus"
	"html/template"
	"io"
	"net/http"
	"net/url"
)

type User struct {
	UserID   int
	Username string `json:"username"`
	Password string `json:"password"`
}

type Profile struct {
	Username        string `json:"username"`
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	Dateofbirthyear string `json:"datebirth"`
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
	logrus.Info("edit_profile opened")
	switch r.Method {
	case http.MethodGet:
		username := r.URL.Query().Get("username")
		if username == "" {
			http.Error(w, "Username not provided", http.StatusBadRequest)
			return
		}
		logrus.Info("Username is:", username)
		ts, err := template.ParseFiles("ui/pages/edit_profile.html")
		if err != nil {
			logrus.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = ts.Execute(w, nil)
		if err != nil {
			logrus.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	case http.MethodPost:
		logrus.Info("post method used")
		edit(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func edit(w http.ResponseWriter, r *http.Request) {
	logrus.Info("edit func opened")
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form data: %s", err), http.StatusBadRequest)
		return
	}
	username := r.URL.Query().Get("username")
	logrus.Info("Requested URL:", r.URL.String())

	data := Profile{
		Username:        username,
		Name:            r.FormValue("name"),
		Surname:         r.FormValue("surname"),
		Dateofbirthyear: r.FormValue("datebirth"),
	}

	// Perform data validation here
	logrus.Info(data)
	err := insertdata(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating user in database: %s", err), http.StatusInternalServerError)
		return
	}
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
		ts, err := template.ParseFiles("ui/pages/Login.html")
		if err != nil {
			logrus.Error(err.Error())
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

		logrus.Info("Received login request for username:", user.Username)
		auth := authenticateUser(user.Username, user.Password)
		exists := userExists(user.Username)
		logrus.Infof("Authentication: %v, Exists: %v for username: %s\n", auth, exists, user.Username)
		if auth && exists {
			logrus.Info("Login successful for username:", user.Username)
			http.Redirect(w, r, "/edit_profile?username="+url.QueryEscape(user.Username), http.StatusSeeOther)
			return
		} else {
			logrus.Info("Login failed or user does not exist for username:", user.Username)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
		}

	}
}
