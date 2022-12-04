package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

func main() {
	os.Remove("./products.db")

	db, err := sql.Open("sqlite3", "products.db")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `CREATE TABLE products_table (ID INTEGER NOT NULL PRIMARY KEY, name TEXT, description TEXT, price REAL);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	var products []product
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
		_, err = stmt.Exec(prod.Id, prod.Name, prod.Description, prod.Price)
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

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.Price); err != nil {
			log.Fatal(err)
		}
		products = append(products, p)
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
	log.Println(products)
	return

}
