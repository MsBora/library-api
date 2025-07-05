package main

import (
	"fmt"
	"library-api/handlers"
	"library-api/storage"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, reading from system environments")
	}

	// configuration
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "mysecretpassword")
	dbName := getEnv("DB_NAME", "postgres")

//--------------

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost,
			dbPort,
			dbUser,
			dbPassword,
			dbName,
	)

	fmt.Println("-----Initializing Storage --------")
	var store *storage.Storage
	err = nil
	maxRetries, i := 5, 0
	for i < maxRetries {
		store, err = storage.NewStorage(connStr)

		if err == nil {
			break
		}

		log.Printf("failed to connect database, retrying in 5 seconds... (%d/%d)", i+1, maxRetries)
		time.Sleep(5 * time.Second)

		i++
	}

	if err != nil {
		log.Fatalf("failed to initialize database after %d retries: %v", i-1, err)
	}
	

	defer store.Close()
	log.Println("storage initialized successfully.")

	bookHandler := handlers.NewBookHandler(store)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/books", func(r chi.Router) {
		r.Post("/", bookHandler.HandleCreateBook)
		r.Get("/", bookHandler.HandleGetBooks)
		r.Get("/{id}", bookHandler.HandleGetBookbyId)
		r.Put("/{id}", bookHandler.HandleUpdateBook)
		r.Delete("/{id}", bookHandler.HandleDeleteBook)
	})

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("failed to start server: %v", err)
	}
}