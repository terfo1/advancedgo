package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"html/template"
	"net/http"
	"os"
)

func main() {
	var limiter = rate.NewLimiter(1, 3)
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	file, errorr := os.OpenFile("log.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if errorr == nil {
		log.Out = file
	} else {
		log.Info("Failed to log into file")
	}
	db := connectDB()
	defer db.Close()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		ts, err := template.ParseFiles("ui/pages/index.html")
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = ts.Execute(w, nil)
	})
	http.HandleFunc("/register", registrationHandler)
	http.HandleFunc("/users", getUsersHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/edit_profile", edit_Profile)
	http.HandleFunc("/listjobs", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("ui/pages/Jobs.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		data := map[string]interface{}{}
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		GetJobs(db, w, r)
	})
	port := 8080
	log.Infof("Server is running on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
	log.WithFields(logrus.Fields{
		"action": "start",
		"status": "success",
	}).Info("App started successfully")
}
