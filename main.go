package main

import (
	"database/sql"
	"log"
	"net/http"

	"shopping-app/pkg/handlers"
	"shopping-app/pkg/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// var tmpl *template.Template
var db *sql.DB

var Store = sessions.NewCookieStore([]byte("shoppingcart"))

func init() {
	//tmpl, _ = template.ParseGlob("templates/*.html")

	/* templatesDir := "./templates"
	pattern := filepath.Join(templatesDir, "**", "*.html")
	tmpl = template.Must(template.ParseGlob(pattern)) */

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
	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	/* fileServer := http.FileServer(http.Dir("./uploads"))
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads", fileServer)) */

	repo := repository.NewRepository(db)
	handler := handlers.NewHandler(repo)

	//User Views and Pages
	r.HandleFunc("/", handler.ShoppingHomepage).Methods("GET")
	r.HandleFunc("/shoppingitems", handler.ShoppingItemsView).Methods("GET")
	r.HandleFunc("/cartitems", handler.CartView).Methods("GET")
	r.HandleFunc("/addtocart/{product_id}", handler.AddToCart).Methods("POST")
	r.HandleFunc("/gotocart", handler.ShoppingCartView).Methods("GET")
	r.HandleFunc("/updateorderitem", handler.UpdateOrderItemQuantity).Methods("PUT")
	r.HandleFunc("/ordercomplete", handler.PlaceOrder).Methods("GET")

	// Product routes
	//Admin Views and Pages
	r.HandleFunc("/manageproducts", handler.ProductsPage).Methods("GET")
	r.HandleFunc("/createproduct", handler.CreateProductView).Methods("GET")
	r.HandleFunc("/allproducts", handler.AllProductsView).Methods("GET")
	r.HandleFunc("/editproduct/{id}", handler.EditProductView).Methods("GET")

	//Admin Actions
	r.HandleFunc("/products", handler.ListProducts).Methods("GET")
	r.HandleFunc("/products", handler.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", handler.GetProduct).Methods("GET")
	r.HandleFunc("/products/{id}", handler.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", handler.DeleteProduct).Methods("DELETE")
	r.HandleFunc("/seed-products", handler.SeedProducts).Methods("POST")

	// Order routes
	r.HandleFunc("/manageorders", handler.OrdersPage).Methods("GET")
	r.HandleFunc("/allorders", handler.AllOrdersView).Methods("GET")
	r.HandleFunc("/orders", handler.ListOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", handler.GetOrder).Methods("GET")

	r.HandleFunc("/orders", handler.CreateOrder).Methods("POST")
	r.HandleFunc("/orders/{id}/items", handler.AddOrderItem).Methods("POST")

	http.ListenAndServe(":5000", r)
}
