package main

import (
	"database/sql"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	// hostname used by the server
	hostname = envOrDefault("HTTP_HOST", "localhost")

	// port used by the server
	port = envOrDefault("PORT", "3000")
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("couldn't loading .env file: %w", err))
	}

	// Connect to the PostgreSQL instance
	db, err := sql.Open("postgres", (&url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD")),
		Host:   net.JoinHostPort(os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT")),
		Path:   os.Getenv("POSTGRES_DB"),
	}).String())
	if err != nil {
		panic(fmt.Errorf("failed ot connect to PostgreSQL: %w", err))
	}
	defer db.Close()

	// Query the Qdrant instance
	qdrantBaseURL := (&url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(os.Getenv("QDRANT_HOST"), os.Getenv("QDRANT_PORT")),
	}).String()
	resp, err := http.Get(qdrantBaseURL + "/v1/collections")
	if err != nil {
		panic(fmt.Errorf("failed to query Qdrant DB: %w", err))
	}
	defer resp.Body.Close()

	// Declare your http routes

	r := chi.NewRouter()

	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "Hello, World!")
	}))

	http.HandleFunc("/livez", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Start the server, which will listen on http://localhost:8080
	addr := net.JoinHostPort("localhost", "8080")
	fmt.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(fmt.Errorf("failed to start server on %s: %w", addr, err))
	}
}

// envOrDefault returns the value of the environment variable at the given key.
// Fallbacks to the given default if the value found is missing or empty.
func envOrDefault(key string, defaultValue string) string {
	if v := strings.TrimSpace(os.Getenv(key)); len(v) > 0 {
		return v
	}
	return defaultValue
}
