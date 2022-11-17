package main

import (
	"github.com/gilmored/pivottech/marvel"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	pubKey := os.Getenv("MARVEL_PUBLIC_KEY")
	privKey := os.Getenv("MARVEL_PRIVATE_KEY")

	client := marvel.MarvelClient{
		BaseURL: "https://gateway.marvel.com:443/v1/public",
		PubKey:  pubKey,
		PrivKey: privKey,
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	characters, err := client.GetCharacters(4)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range characters {
		log.Println("character:", v)
	}

}
