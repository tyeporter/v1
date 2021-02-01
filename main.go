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
)

func main() {
	// Create a new router
	router := router.Router()

	// Get port
	port := os.Getenv("PORT")

	// Start server...
	fmt.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
