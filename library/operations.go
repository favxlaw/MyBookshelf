package library

import (
	"fmt"
	"github.com/favxlaw/models"
	"time"
)

// AddBook adds a new book to the library
func AddBook(library []models.Book, book models.Book) []models.Book {
	return append(library, book)
}

// ListAllBooks displays all books in the library
func ListAllBooks(library []models.Book) {
	if len(library) == 0 {
		fmt.Println("No books in library")
		return
	}

	fmt.Println("\n Your Library:")
	fmt.Println("================")
	for i, book := range library {
		fmt.Printf("%d. %s by %s [%s]\n", i+1, book.Title, book.Author, book.Status)
	}
}

// GetBookByTitle finds and returns a book by its title
func GetBookByTitle(library []models.Book, title string) *models.Book {
	for i := range library {
		if library[i].Title == title {
			return &library[i]
		}
	}
	return nil
}

// UpdateBookStatus changes the status of a book and sets EndDate if finished
func UpdateBookStatus(library []models.Book, title string, newStatus models.BookStatus) bool {
	book := GetBookByTitle(library, title)
	if book == nil {
		return false
	}

	book.Status = newStatus

	// If marking as finished, set the end date
	if newStatus == models.StatusFinished || newStatus == models.StatusAbandoned {
		now := time.Now()
		book.EndDate = &now
	}

	return true
}

// DeleteBook removes a book from the library by title
func DeleteBook(library []models.Book, title string) []models.Book {
	for i, book := range library {
		if book.Title == title {
			return append(library[:i], library[i+1:]...)
		}
	}
	return library
}

// FilterByStatus returns all books with a specific status
func FilterByStatus(library []models.Book, status models.BookStatus) []models.Book {
	var result []models.Book
	for _, book := range library {
		if book.Status == status {
			result = append(result, book)
		}
	}
	return result
}

// FilterByCategory returns all books in a specific category
func FilterByCategory(library []models.Book, category string) []models.Book {
	var result []models.Book
	for _, book := range library {
		if book.Category == category {
			result = append(result, book)
		}
	}
	return result
}
