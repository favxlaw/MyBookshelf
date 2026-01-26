package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/favxlaw/config"
	"github.com/favxlaw/handlers"
	"github.com/favxlaw/models"
	"github.com/favxlaw/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	log.Printf("Starting Book Tracker API")
	log.Printf("Configuration:")
	log.Printf("  Port: %s", cfg.Port)
	log.Printf("  Database: %s", cfg.DBPath)
	log.Printf("  Log Level: %s", cfg.LogLevel)
	log.Println()

	bookStore, err := store.NewSQLiteStore(cfg.DBPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer bookStore.Close()

	seedBooks(bookStore)

	bookHandler := handlers.NewBookHandler(bookStore)

	http.Handle("/books", bookHandler)
	http.Handle("/books/", bookHandler)
	http.HandleFunc("/", homeHandler)

	fmt.Println("Server starting on http://localhost:" + cfg.Port)
	fmt.Println("GET    /books       - List all books")
	fmt.Println("POST   /books       - Add new book")
	fmt.Println("GET    /books/{id}  - Get specific book")
	fmt.Println("PUT    /books/{id}  - Update book")
	fmt.Println("DELETE /books/{id}  - Delete book")
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop")

	port := cfg.Port
	if port[0] != ':' {
		port = ":" + port
	}

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func seedBooks(s *store.SQLiteStore) {
	// Only seed if database is empty
	existing := s.GetAll()
	if len(existing) > 0 {
		return
	}

	fmt.Println("Seeding initial data...")

	s.Create(models.Book{
		Title:     "Clean Code",
		Author:    "Robert C. Martin",
		Status:    models.StatusToRead,
		Category:  "Software Engineering",
		StartDate: time.Now(),
	})

	s.Create(models.Book{
		Title:     "Dune",
		Author:    "Frank Herbert",
		Status:    models.StatusReading,
		Category:  "Science Fiction",
		StartDate: time.Now(),
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Book Tracker API\n\n")
	fmt.Fprintf(w, "Available Endpoints:\n")
	fmt.Fprintf(w, "  GET    /books       - List all books\n")
	fmt.Fprintf(w, "  POST   /books       - Add new book\n")
	fmt.Fprintf(w, "  GET    /books/{id}  - Get specific book\n")
	fmt.Fprintf(w, "  PUT    /books/{id}  - Update book\n")
	fmt.Fprintf(w, "  DELETE /books/{id}  - Delete book\n")
}
