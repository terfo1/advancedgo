package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

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
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/edit_profile", edit_Profile)
	port := 8080
	fmt.Printf("Server is running on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
