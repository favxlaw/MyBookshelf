package main

import (
	"encoding/json"
	"fmt"
	"github.com/favxlaw/library"
	"github.com/favxlaw/models"
	"net/http"
	"time"
)

// in-memory library (for now)
var myBook []models.Book

func main() {
	myBook = library.AddBook(myBook, models.Book{
		Title:     "The Pragmatic Programmer",
		Author:    "Hunt & Thomas",
		Status:    models.StatusReading,
		Category:  "Software Engineering",
		StartDate: time.Now(),
	})

	myBook = library.AddBook(myBook, models.Book{
		Title:     "Clean Code",
		Author:    "Robert C. Martin",
		Status:    models.StatusToRead,
		Category:  "Software Engineering",
		StartDate: time.Now(),
	})

	myBook = library.AddBook(myBook, models.Book{
		Title:     "Dune",
		Author:    "Frank Herbert",
		Status:    models.StatusReading,
		Category:  "Science Fiction",
		StartDate: time.Now(),
	})

	// Register routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/books", booksHandler)

	// Start server
	fmt.Println("Server starting on http://localhost:8006")
	fmt.Println("Try: http://localhost:8006/books")
	http.ListenAndServe(":8006", nil)
}

// homeHandler handles requests to "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Book Tracker API\n\nEndpoints:\nGET /books - List all books")
}

// booksHandler handles requests to "/books"
func booksHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Convert books to JSON
	json.NewEncoder(w).Encode(myBook)
}
