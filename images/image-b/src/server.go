package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"structs"

	"registerService"

	"github.com/joho/godotenv"
)

/*
Load dev environment from .env if os.Getenv('GO_ENV') != 'prod'
*/
func loadDevEnv() {
	// Get dev env if not prod
	if os.Getenv("GO_ENV") != "prod" {
		// Load .env file
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

}

/*
Registration endpoint
Accepts a POST request containing RegisterInfo as JSON.
Returns:
- 400 if request is invalid (bad content type, bad method or invalid JSON)
-
*/
func registerEndpoint(w http.ResponseWriter, req *http.Request) {
	// Validate request is a POST w/ JSON
	if req.Header.Get("Content-Type") != "application/json" || req.Method != "POST" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	// Get RegisterInfo from JSON body
	var body structs.RegisterInfo
	err := json.NewDecoder(req.Body).Decode(&body)
	// Fail if json was invalid
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// Handle business logic
	res, resCode := registerService.RegisterUser()
}

func main() {
	loadDevEnv()
	http.HandleFunc("/register", registerEndpoint)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
