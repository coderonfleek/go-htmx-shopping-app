package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"shopping-app/pkg/handlers"
	"shopping-app/pkg/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var tmpl *template.Template
var db *sql.DB

var Store = sessions.NewCookieStore([]byte("usermanagementsecret"))

func init() {
	tmpl, _ = template.ParseGlob("templates/*.html")

	//Set up Sessions
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 3,
		HttpOnly: true,
	}

}

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

	// Setup Static file handling for images

	fileServer := http.FileServer(http.Dir("./uploads"))
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads", fileServer))

	repo := repository.NewRepository(db)
	handler := handlers.NewHandler(repo)

	// Product routes
	r.HandleFunc("/products", handler.ListProducts).Methods("GET")
	r.HandleFunc("/products", handler.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", handler.GetProduct).Methods("GET")
	r.HandleFunc("/products/{id}", handler.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", handler.DeleteProduct).Methods("DELETE")

	// Order routes
	r.HandleFunc("/orders", handler.CreateOrder).Methods("POST")
	r.HandleFunc("/orders/{id}", handler.GetOrder).Methods("GET")
	r.HandleFunc("/orders/{id}/items", handler.AddOrderItem).Methods("POST")

	http.ListenAndServe(":5000", r)
}
