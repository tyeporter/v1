/*
	The router package defines all of the API endpoints.
*/

package router

import (
	"backend/middleware"

	"github.com/gorilla/mux"
)

// Router returns a new router of type [*mux.Router].
// It creates all the endpoints and the respective middleware.
func Router() *mux.Router {
	// Create new router
	router := mux.NewRouter()

	// fs := http.FileServer(http.Dir("./build/"))
	// router.PathPrefix("/").Handler(fs)

	router.HandleFunc("/api/articles/{name}", middleware.GetArticle).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/articles/{name}/like", middleware.LikeArticle).Methods("POST", "OPTIONS")
	spa := middleware.SPAHandler{StaticPath: "build", IndexPath: "index.html"}
	router.PathPrefix("/").Handler(spa).Methods("GET", "OPTIONS")

	return router
}
