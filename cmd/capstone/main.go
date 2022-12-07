package main

import "github.com/gilmored/pivottech/capstone"

//import "github.com/gilmored/pivottech/capstone"

//import (
//	"encoding/json"
//	"fmt"
//	"github.com/joho/godotenv"
//	"log"
//	"net/http"
//	"os"
//	"time"
//)
//
//type veh struct {
//	Data struct {
//		Year         int    `json:"year"`
//		Make         string `json:"make"`
//		Model        string `json:"model"`
//		Manufacturer string `json:"manufacturer"`
//		Engine       string `json:"engine"`
//		Trim         string `json:"trim"`
//		Transmission string `json:"transmission"`
//	}
//}
//
//func getVehicleInfo() {
//	err := godotenv.Load()
//	if err != nil {
//		log.Println("Error loading .env file")
//	}
//	ptnrTok := os.Getenv("CARMD_PARTNER_TOKEN")
//	authKey := os.Getenv("CARMD_AUTH_KEY")
//	var vin string
//	fmt.Println("Enter a VIN to see information about your vehicle. If you don't have one enter the following for testing: 1GNALDEK9FZ108495")
//	_, err = fmt.Scanln(&vin)
//	url := "https://api.carmd.com/v3.0/decode?vin=" + vin
//	c := http.Client{Timeout: time.Duration(2) * time.Second}
//	//req, err := http.NewRequest("GET", "https://api.carmd.com/v3.0/decode?vin=1GNALDEK9FZ108495", nil)
//	//req, err := http.NewRequest("GET", "https://api.carmd.com/v3.0/decode?vin=5uxwx5c54bl703157", nil)
//	req, err := http.NewRequest("GET", url, nil)
//	fmt.Println(url)
//
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	req.Header.Add("Content-Type", "application/json")
//	req.Header.Add("authorization", "Basic "+authKey)
//	req.Header.Add("partner-token", ptnrTok)
//
//	resp, err := c.Do(req)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	defer resp.Body.Close()
//
//	var veh veh
//
//	if err := json.NewDecoder(resp.Body).Decode(&veh); err != nil {
//		log.Println("Cannot locate VIN. Please try again.")
//		return
//	}
//	if veh.Data.Year == 0 {
//		log.Println("You have no api credits remaining to view the returning information.  Please try again tomorrow.")
//		log.Println("Response Status:", resp.Status)
//	} else {
//		log.Println("Response Status:", resp.Status)
//		log.Printf("\nYear: %v\n Make: %v\n Model: %v\n Engine: %v\n Trim: %v\n Transmission: %v\n", veh.Data.Year, veh.Data.Make, veh.Data.Model, veh.Data.Engine, veh.Data.Trim, veh.Data.Transmission)
//		fmt.Println("Click the link below to see a catalog of replacement parts for your vehicle. Please note, you may need to reference the engine information above to make a selection in the catalog.")
//		log.Printf("https://www.rockauto.com/en/catalog/%v,%v,%v", veh.Data.Make, veh.Data.Year, veh.Data.Model)
//	}
//veh.Data.Make = "CHEVROLET"
//veh.Data.Year = 2015
//veh.Data.Model = "EQUINOX"
//log.Printf("\nYear: %v\n Make: %v\n Model: %v\n Engine: %v\n Trim: %v\n Transmission: %v\n", veh.Data.Year, veh.Data.Make, veh.Data.Model, veh.Data.Engine, veh.Data.Trim, veh.Data.Transmission)
//log.Println("Response Status:", resp.Status)
//log.Printf("https://www.rockauto.com/en/catalog/%v,%v,%v", veh.Data.Make, veh.Data.Year, veh.Data.Model)

//}

//func main() {
//
//	getVehicleInfo()
//}

func main() {
	capstone.GetVehicleInfo(capstone.GetVin())
}
