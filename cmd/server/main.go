package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)

type product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

var products []product

//type prodDB struct {
//	id          int
//	name        string
//	description string
//	price       int64
//}
//
//var prods []prodDB

var db *sql.DB
var err error

func connectDB() *sql.DB {
	db, err = sql.Open("sqlite3", "/Users/dustin/go/src/github.com/gilmored/pivottech/cmd/seeder/products.db")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	rows, err := db.Query("SELECT id, name, description, price FROM products_table")
	if err != nil {
		log.Fatal(err)
	}
	//defer rows.Close()

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

	//err = rows.Err()
	//if err != nil {
	//	log.Fatal(err)
	//}

	status := "status: up"
	if err := db.Ping(); err != nil {
		status = "status: down"
		log.Fatalf("%s error: %v", status, err)
	}
	log.Println(status)
	log.Println(products)
	return db
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	///
	//vars := mux.Vars(r)
	//limit, err := strconv.ParseInt(vars["limit"], 10, 64)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//if limit < 0 || limit > int64(len(products)) {
	//	total := len(products)
	//	log.Printf("Selected limit is not within Products range. Please select a limit between 1 and %d", total)
	//} else {
	//	if _, err = db.Query("SELECT * FROM products_table LIMIT ?", limit); err != nil {
	//		log.Println("error", err)
	//		w.WriteHeader(http.StatusInternalServerError)
	//	}
	//}
	if _, err = db.Query("SELECT * FROM products_table"); err != nil {
		log.Println("error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//func getProducts(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	if err := json.NewEncoder(w).Encode(products); err != nil {
//		log.Println("error", err)
//		w.WriteHeader(http.StatusInternalServerError)
//	}
//}

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
	connectDB()
	r := mux.NewRouter()
	//r.HandleFunc("/api/products?limit={limit}", getProducts).Methods("GET")
	//
	r.HandleFunc("/api/products", getProducts).Methods("GET")
	//r.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	//r.HandleFunc("/api/products", createProduct).Methods("POST")
	//r.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	//r.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")
	//
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
