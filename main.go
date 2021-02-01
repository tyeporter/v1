/*
	The main package will start the server.
*/

package main

import (
	"backend/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Create a new router
	router := router.Router()

	// Start server...
	fmt.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
