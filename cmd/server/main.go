package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

type product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

var products []product

var db *sql.DB
var err error

func connectDB() *sql.DB {
	dbPath, err := filepath.Abs("../../cmd/seeder/products.db")
	if err != nil {
		log.Println(err)
	}
	db, err = sql.Open("sqlite3", dbPath)

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id, name, description, price FROM products_table")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var description string
		var price float64
		if err := rows.Scan(&id, &name, &description, &price); err != nil {
			log.Fatal(err)
		}
		products = append(products, product{ID: id, Name: name, Description: description, Price: price})
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	status := "status: up"
	if err := db.Ping(); err != nil {
		status = "status: down"
		log.Fatalf("%s error: %v", status, err)
	}
	log.Println(status)
	return db
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	limit, err := strconv.ParseInt(params["limit"], 10, 64)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if limit < 0 || limit > int64(len(products)) {
		w.WriteHeader(400)
		return
	}

	rows, err := db.Query("SELECT id, name, description, price FROM products_table LIMIT ?", limit)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var description string
		var price float64
		if err := rows.Scan(&id, &name, &description, &price); err != nil {
			log.Fatal(err)
		}
		products = append(products, product{ID: id, Name: name, Description: description, Price: price})
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(w).Encode(products); err != nil {
		w.WriteHeader(500)
		return
	}
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	paramsConvert, err := strconv.Atoi(params["id"])
	if paramsConvert < 1 {
		w.WriteHeader(404)
		return
	}
	if err != nil {
		w.WriteHeader(400)
		return
	}

	var products product

	rows, err := db.Query("SELECT id, name, description, price FROM products_table WHERE id = ?", params["id"])
	if err != nil {
		w.WriteHeader(400)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var description string
		var price float64
		if err := rows.Scan(&id, &name, &description, &price); err != nil {
			log.Fatal(err)
		}
		products = product{
			ID:          id,
			Name:        name,
			Description: description,
			Price:       price,
		}
	}

	if products.ID == 0 {
		w.WriteHeader(400)
		return
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(w).Encode(products); err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var prod product
	if err := json.NewDecoder(r.Body).Decode(&prod); err != nil {
		w.WriteHeader(400)
		return
	}
	if prod.Name == "" {
		w.WriteHeader(400)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	stmt, err := tx.Prepare("INSERT INTO products_table (name, description, price) VALUES (?, ?, ?)")
	if err != nil {
		w.WriteHeader(400)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(prod.Name, prod.Description, prod.Price)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	if err := tx.Commit(); err != nil {
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(201)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	paramsConvert, err := strconv.Atoi(params["id"])
	if paramsConvert < 0 {
		w.WriteHeader(400)
		return
	}
	if err != nil {
		w.WriteHeader(400)
		return
	}

	var products product
	if err := json.NewDecoder(r.Body).Decode(&products); err != nil {
		w.WriteHeader(400)
		return
	}
	if products.Name == "" || products.Price == 0 {
		w.WriteHeader(400)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	stmt, err := tx.Prepare("UPDATE products_table SET name = ?, description = ?, price = ? WHERE id = ?")
	if err != nil {
		w.WriteHeader(400)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(products.Name, products.Description, products.Price, params["id"])
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if err := tx.Commit(); err != nil {
		w.WriteHeader(400)
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if rows == 0 {
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(200)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	if _, err := strconv.Atoi(params["id"]); err != nil {
		w.WriteHeader(400)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		w.WriteHeader(400)
		return
	}

	stmt, err := tx.Prepare("DELETE FROM products_table WHERE id = ?")
	if err != nil {
		w.WriteHeader(400)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(params["id"])
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if err := tx.Commit(); err != nil {
		w.WriteHeader(400)
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if rows == 0 {
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(200)

}

func main() {
	connectDB()
	r := mux.NewRouter()

	r.HandleFunc("/api/products/limit={limit}", getProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	r.HandleFunc("/api/products", createProduct).Methods("POST")
	r.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
