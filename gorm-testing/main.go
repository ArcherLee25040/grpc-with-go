package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

var db *gorm.DB

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Product{})
}

func main() {
	initDB()
	// Seed the database with a product if necessary
	db.FirstOrCreate(&Product{Code: "D42", Price: 100})

	http.HandleFunc("/", homeHandler) // Add a handler for the root path
	http.HandleFunc("/products", getProducts)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Product API! Visit /products to see all products.")
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	var products []Product
	db.Find(&products)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
