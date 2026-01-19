// package main

// import (
// 	"fmt"
// 	"github.com/favxlaw/library"
// 	"github.com/favxlaw/models"
// 	"time"
// )

// //func main() {
// 	// Initialize empty library
// 	myLibrary := []models.Book{}

// 	// Add books
// 	myLibrary = library.AddBook(myLibrary, models.Book{
// 		Title:     "The Pragmatic Programmer",
// 		Author:    "Hunt & Thomas",
// 		Status:    models.StatusReading,
// 		Category:  "Software Engineering",
// 		StartDate: time.Now(),
// 	})

// 	myLibrary = library.AddBook(myLibrary, models.Book{
// 		Title:     "Clean Code",
// 		Author:    "Robert C. Martin",
// 		Status:    models.StatusToRead,
// 		Category:  "Software Engineering",
// 		StartDate: time.Now(),
// 	})

// 	myLibrary = library.AddBook(myLibrary, models.Book{
// 		Title:     "Dune",
// 		Author:    "Frank Herbert",
// 		Status:    models.StatusReading,
// 		Category:  "Science Fiction",
// 		StartDate: time.Now(),
// 	})

// 	// Display all books
// 	library.ListAllBooks(myLibrary)

// 	// Update a book status
// 	fmt.Println("\n Updating 'The Pragmatic Programmer' to finished...")
// 	library.UpdateBookStatus(myLibrary, "The Pragmatic Programmer", models.StatusFinished)
// 	library.ListAllBooks(myLibrary)

// 	// Delete a book
// 	fmt.Println("\n  Deleting 'Clean Code'...")
// 	myLibrary = library.DeleteBook(myLibrary, "Clean Code")
// 	library.ListAllBooks(myLibrary)

// 	// Filter by status
// 	fmt.Println("\n Currently Reading:")
// 	reading := library.FilterByStatus(myLibrary, models.StatusReading)
// 	for _, book := range reading {
// 		fmt.Printf("- %s\n", book.Title)
// 	}

// 	// Filter by category
// 	fmt.Println("\n Software Engineering Books:")
// 	techBooks := library.FilterByCategory(myLibrary, "Software Engineering")
// 	for _, book := range techBooks {
// 		fmt.Printf("- %s by %s\n", book.Title, book.Author)
// 	}
// }

package main

import (
	"fmt"
	"github.com/favxlaw/handlers"
	"github.com/favxlaw/models"
	"github.com/favxlaw/store"
	"log"
	"net/http"
	"time"
)

func main() {
	bookStore, err := store.NewSQLiteStore("./booktracker.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer bookStore.Close()

	// Seed with initial data that's only if the database is empty
	seedBooks(bookStore)

	// Create handler
	bookHandler := handlers.NewBookHandler(bookStore)

	// Register routes
	http.Handle("/books", bookHandler)
	http.Handle("/books/", bookHandler)
	http.HandleFunc("/", homeHandler)

	// Start server
	fmt.Println("Server starting on http://localhost:8006")
	fmt.Println("GET    /books       - List all books")
	fmt.Println("POST   /books       - Add new book")
	fmt.Println("GET    /books/{id}  - Get specific book")
	fmt.Println("PUT    /books/{id}  - Update book")
	fmt.Println("DELETE /books/{id}  - Delete book")
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop")

	http.ListenAndServe(":8006", nil)
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
