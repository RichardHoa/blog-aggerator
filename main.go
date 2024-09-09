package main

import (
	// "fmt"
	"database/sql"
	"github.com/RichardHoa/blog-aggerator/internal/config"
	"github.com/RichardHoa/blog-aggerator/internal/database"
	"github.com/RichardHoa/blog-aggerator/internal/route"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {

	const (
		fetchInterval = 60 * time.Second
		numFeeds      = 10
	)

	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	PORT := os.Getenv("PORT")
	dbURL := os.Getenv("PSQL_CONNECTION_STRING")

	// Open connection to SQL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Create new db queries
	dbQueries := database.New(db)

	// Add db queries to api config
	apiCfg := &config.ApiConfig{
		DB: dbQueries,
	}

	// Create new server
	mux := http.NewServeMux()

	// Configure routes
	route.ConfigureRoutes(mux, apiCfg)

	// Create server
	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: mux,
	}

	stop := make(chan struct{})
	go feedWorker(stop, fetchInterval, numFeeds, apiCfg)

	// Run server
	log.Printf("Serving files from %s on port: %s\n", ".", PORT)
	log.Fatal(server.ListenAndServe())

}
