package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
)

const (
	user     = "postgres"
	password = "terfo2005"
	dbname   = "postgres"
)

type Job struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Company     string `json:"company"`
	Description string `json:"description"`
	AddedDate   string `json:"added_date"`
}

func connectDB() *sql.DB {
	psqlInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func GetJobs(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	nameFilter := r.URL.Query().Get("name")
	companyFilter := r.URL.Query().Get("company")
	sort := r.URL.Query().Get("sort")
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")

	// Default values for pagination
	if page == "" {
		page = "1"
	}
	if pageSize == "" {
		pageSize = "3"
	}

	pageNum, _ := strconv.Atoi(page)
	pageSizeNum, _ := strconv.Atoi(pageSize)
	offset := (pageNum - 1) * pageSizeNum

	query := `SELECT id, name, company, description, added_date FROM jobs WHERE 1=1`

	// Filtering
	args := []interface{}{}
	if nameFilter != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", len(args)+1)
		args = append(args, "%"+nameFilter+"%")
	}
	if companyFilter != "" {
		query += fmt.Sprintf(" AND company ILIKE $%d", len(args)+1)
		args = append(args, "%"+companyFilter+"%")
	}

	// Sorting
	if sort != "" {
		query += fmt.Sprintf(" ORDER BY %s", sort)
	} else {
		query += " ORDER BY added_date DESC"
	}

	// Pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, pageSizeNum, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer rows.Close()

	jobs := make([]Job, 0)
	for rows.Next() {
		var job Job
		if err := rows.Scan(&job.ID, &job.Name, &job.Company, &job.Description, &job.AddedDate); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		jobs = append(jobs, job)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}
