package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to the PostgreSQL instance
	db, err := sql.Open("postgres", (&url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD")),
		Host:   net.JoinHostPort(os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT")),
		Path:   os.Getenv("POSTGRES_DB"),
	}).String())
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Query the Qdrant instance
	qdrantBaseURL := (&url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(os.Getenv("QDRANT_HOST"), os.Getenv("QDRANT_PORT")),
	}).String()
	resp, err := http.Get(qdrantBaseURL + "/v1/collections")
	if err != nil {
		log.Fatalf("Error connecting to Qdrant: %v", err)
	}
	defer resp.Body.Close()

	// Declare your http routes

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Hello, World!")
	})

	http.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Start the server, which will listen on http://localhost:8080
	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
