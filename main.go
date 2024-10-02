package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func initDB() {
	var err error
	// Initialize the db variable
	db, err = sql.Open("mysql", "root:root@(127.0.0.1:3333)/shopping?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	// Check the database connection
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	r := mux.NewRouter()

	//Setup MySQL
	initDB()
	defer db.Close()

	// Setup Static folder for static files and images
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":5000", r)
}
