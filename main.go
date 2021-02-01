/*
	The main package will start the server.
*/

package main

import (
	"backend/router"
	"fmt"
	"log"
	"net/http"
	"os"
	// "github.com/joho/godotenv" // FOR DEVELOPMENT PURPOSES
)

func main() {
	// FOR DEVELOPMENT PURPOSES
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	// Create a new router
	router := router.Router()

	// Get port
	port := os.Getenv("PORT")

	// Start server...
	fmt.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
