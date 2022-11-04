package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

var products []product

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func initProducts() {
	bs, err := os.ReadFile("products.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(bs, &products); err != nil {
		log.Fatal(err)
	}
}

func main() {
	initProducts()
	r := mux.NewRouter()
	r.HandleFunc("/products", getProductsHandler)

	r.HandleFunc("/api/products", getProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	r.HandleFunc("/api/products", createProduct).Methods("POST")
	r.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)

}

func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range products {
		//if string(item.ID) == params["id"] {
		if strconv.Itoa(item.ID) == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var prod product
	_ = json.NewDecoder(r.Body).Decode(&prod)
	prod.ID = rand.Intn(int(time.Now().UnixNano()))
	products = append(products, prod)
	json.NewEncoder(w).Encode(prod)

}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var prod product
	for index, item := range products {
		if strconv.Itoa(item.ID) == params["id"] {
			products = append(products[:index], products[index+1:]...)
			_ = json.NewDecoder(r.Body).Decode(&prod)
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
	for index, item := range products {
		if strconv.Itoa(item.ID) == params["id"] {
			products = append(products[:index], products[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(products)

}
