package main

import (
	"fmt"
	"github.com/favxlaw/library"
	"github.com/favxlaw/models"
	"time"
)

func main() {
	// Initialize empty library
	myLibrary := []models.Book{}

	// Add books
	myLibrary = library.AddBook(myLibrary, models.Book{
		Title:     "The Pragmatic Programmer",
		Author:    "Hunt & Thomas",
		Status:    models.StatusReading,
		Category:  "Software Engineering",
		StartDate: time.Now(),
	})

	myLibrary = library.AddBook(myLibrary, models.Book{
		Title:     "Clean Code",
		Author:    "Robert C. Martin",
		Status:    models.StatusToRead,
		Category:  "Software Engineering",
		StartDate: time.Now(),
	})

	myLibrary = library.AddBook(myLibrary, models.Book{
		Title:     "Dune",
		Author:    "Frank Herbert",
		Status:    models.StatusReading,
		Category:  "Science Fiction",
		StartDate: time.Now(),
	})

	// Display all books
	library.ListAllBooks(myLibrary)

	// Update a book status
	fmt.Println("\n Updating 'The Pragmatic Programmer' to finished...")
	library.UpdateBookStatus(myLibrary, "The Pragmatic Programmer", models.StatusFinished)
	library.ListAllBooks(myLibrary)

	// Delete a book
	fmt.Println("\n  Deleting 'Clean Code'...")
	myLibrary = library.DeleteBook(myLibrary, "Clean Code")
	library.ListAllBooks(myLibrary)

	// Filter by status
	fmt.Println("\n Currently Reading:")
	reading := library.FilterByStatus(myLibrary, models.StatusReading)
	for _, book := range reading {
		fmt.Printf("- %s\n", book.Title)
	}

	// Filter by category
	fmt.Println("\n💻 Software Engineering Books:")
	techBooks := library.FilterByCategory(myLibrary, "Software Engineering")
	for _, book := range techBooks {
		fmt.Printf("- %s by %s\n", book.Title, book.Author)
	}
}
