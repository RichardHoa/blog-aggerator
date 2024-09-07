package main

import (
	// "fmt"
	"log"
	"net/http"
	"os"
	"github.com/RichardHoa/blog-aggerator/internal/route"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	PORT := os.Getenv("PORT")

	mux := http.NewServeMux()

	route.ConfigureRoutes(mux)

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", ".", PORT)
	log.Fatal(server.ListenAndServe())

}
