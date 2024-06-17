package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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

	// Connect to PostgreSQL
	pgConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Connect to Qdrant
	qdrantURL := fmt.Sprintf("http://%s:%s", os.Getenv("QDRANT_HOST"), os.Getenv("QDRANT_PORT"))
	resp, err := http.Get(qdrantURL + "/v1/collections")
	if err != nil {
		log.Fatalf("Error connecting to Qdrant: %v", err)
	}
	defer resp.Body.Close()

	// Simple HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Connected to PostgreSQL and Qdrant successfully!")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
