/*
	The middleware package is the bridge between the API and the database.

	It will handle all the db operations like Insert, Select, Update, and Delete (CRUD).
*/

package middleware

import (
	"backend/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// --------------- HELPERS --------------------

// The connection string to the MongoDB client
var connectionString string

// The name of the database we're connecting to
var databaseName string

// createConnection should be called early in our middleware 
// handler functions that rely on a connection to the MongoDB client
func createConnection() *mongo.Client {
	// Set the connection string and database name
	connectionString = "mongodb://localhost:27017"
	databaseName = "tyeporter-dev"

	// Define 10 second timeout for MongoDB - GO connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		panic(err)
	}

	// Use ping to check that our connection is ok/open
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	// Return MongoDB client
	return client
}

func withDB(operations func(*mongo.Database, context.Context)) {
	// Define 10 second timeout for MongoBD - Go connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create client connection
	client := createConnection()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Connect to database
	database := client.Database(databaseName)
	if database == nil {
		log.Fatalf("Unable to establish a connection with the database (%s)", databaseName)
	}

	operations(database, ctx)
}

// --------------- SPA FUNCTIONS --------------------

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type SPAHandler struct {
	StaticPath string
	IndexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h SPAHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    // get the absolute path to prevent directory traversal
	path, err := filepath.Abs(req.URL.Path)
	if err != nil {
        // if we failed to get the absolute path respond with a 400 bad request
        // and stop
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

    // prepend the path with the path to the static directory
	path = filepath.Join(h.StaticPath, path)

    // check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(res, req, filepath.Join(h.StaticPath, h.IndexPath))
		return
	} else if err != nil {
        // if we got an error (that wasn't that the file doesn't exist) stating the
        // file, return a 500 internal server error and stop
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

    // otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(res, req)
}

// --------------- HANDLER FUNCTIONS --------------------

// GetArticle corresponds to the "/api/articles/{name}" endpoint.
// It retreives a single record from the database.
func GetArticle(res http.ResponseWriter, req *http.Request) {
	// Set header values
	res.Header().Add("Content-Type", "application/json")

	// Get article name from request parameters
	params := mux.Vars(req)
	name := params["name"]
	if name == "" {
		log.Fatal("Unable to get name from request parameters")
	}

	// Create an empty Article
	var article models.Article

	withDB(func (db *mongo.Database, ctx context.Context) {
		// Retrieve the articles collection from database
		collection := db.Collection("articles")

		// Find article by name and marshal into article obejct
		err := collection.FindOne(ctx, bson.M{"name": name}).Decode(&article)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(`{ "message": "` + err.Error() + `"}`))
			return
		}

		// Encode article to response writer
		err = json.NewEncoder(res).Encode(article)
		if err != nil {
			log.Fatal("Unable to encode response")
		}
	})
}

func LikeArticle(res http.ResponseWriter, req *http.Request) {
	// Set header values
	res.Header().Add("Content-Type", "application/json")

	// Get article name from request parameters
	params := mux.Vars(req)
	name := params["name"]
	if name == "" {
		log.Fatal("Unable to get name from request parameters")
	}

	// Create an empty Article
	var article models.Article

	withDB(func (db *mongo.Database, ctx context.Context){
		// Retrieve the articles collection from database
		collection := db.Collection("articles")

		// Find article by name and marshal into article obejct
		err := collection.FindOne(ctx, bson.M{"name": name}).Decode(&article)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(`{ "message": "` + err.Error() + `"}`))
			return
		}

		// Create the update
		update := bson.M{"$set": bson.M{"likes": article.Likes + 1}}

		// Increment the likes on that article
		result := collection.FindOneAndUpdate(ctx, bson.M{"name": article.Name}, update)
		if result.Err() != nil {
			log.Fatal("Unable to update article")
		}

		// Marshal result into new object
		var newVersion models.Article
		result.Decode(&newVersion)

		// newVersion corresponse to the old version so we have to increment the likes on our response
		newVersion.Likes += 1

		// Encode new version of article to response writer
		err = json.NewEncoder(res).Encode(newVersion)
		if err != nil {
			log.Fatal("Unable to encode response")
		}
	})
	
}
