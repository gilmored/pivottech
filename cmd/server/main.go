package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"strconv"
)

type product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

var products []product

func createDB() {
	os.Remove("./products.db")
	db, err := sql.Open("sqlite3", "./products.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
		CREATE TABLE products_table (ID INTEGER NOT NULL PRIMARY KEY, name TEXT, description TEXT, price REAL);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	bs, err := os.ReadFile("./products.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(bs, &products); err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	stmt, err := tx.Prepare("INSERT INTO products_table(id, name, description, price) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	for _, prod := range products {
		_, err = stmt.Exec(prod.ID, prod.Name, prod.Description, prod.Price)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}

	rows, err := db.Query("SELECT id, name, description, price FROM products_table LIMIT 5")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type prodDB struct {
		id          int
		name        string
		description string
		price       int64
	}
	var prods []prodDB
	for rows.Next() {
		var p prodDB
		if err := rows.Scan(&p.id, &p.name, &p.description, &p.price); err != nil {
			log.Fatal(err)
		}
		prods = append(prods, p)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	status := "status: up"
	if err := db.Ping(); err != nil {
		status = "status: down"
	}
	log.Println(status)
	log.Println(prods)
	return

}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Println("error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range products {
		if strconv.Itoa(item.ID) == params["id"] {
			if err := json.NewEncoder(w).Encode(item); err != nil {
				log.Println("error", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var prod product
	if err := json.NewDecoder(r.Body).Decode(&prod); err != nil {
		log.Println("error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	prod.ID = len(products) + 1
	products = append(products, prod)
	json.NewEncoder(w).Encode(prod)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var prod product
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	for index, item := range products {
		if id == item.ID {
			products = append(products[:index], products[index+1:]...)
			if err := json.NewDecoder(r.Body).Decode(&prod); err != nil {
				log.Println("error", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			products = append(products, prod)
			json.NewEncoder(w).Encode(prod)
			return
		}
	}
	json.NewEncoder(w).Encode(products)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	for index, item := range products {
		if id == item.ID {
			products = append(products[:index], products[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(products)

}

func main() {
	createDB()
	r := mux.NewRouter()

	r.HandleFunc("/api/products", getProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	r.HandleFunc("/api/products", createProduct).Methods("POST")
	r.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
