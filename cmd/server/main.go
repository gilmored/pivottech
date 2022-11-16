package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
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

func initProducts() {
	bs, err := os.ReadFile("products.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(bs, &products); err != nil {
		log.Fatal(err)
	}
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
	initProducts()
	r := mux.NewRouter()

	r.HandleFunc("/api/products", getProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	r.HandleFunc("/api/products", createProduct).Methods("POST")
	r.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
