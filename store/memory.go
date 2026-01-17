package store

import (
	"fmt"
	"github.com/favxlaw/models"
)

// BookStore manages book data
type BookStore struct {
	books  []models.Book
	nextID int
}

// NewBookStore creates a new book store
func NewBookStore() *BookStore {
	return &BookStore{
		books:  []models.Book{},
		nextID: 1,
	}
}

// GetAll returns all books
func (s *BookStore) GetAll() []models.Book {
	return s.books
}

// GetByID finds a book by ID
func (s *BookStore) GetByID(id int) (*models.Book, error) {
	for i := range s.books {
		if s.books[i].ID == id {
			return &s.books[i], nil
		}
	}
	return nil, fmt.Errorf("book not found")
}

// Create adds a new book and assigns an ID
func (s *BookStore) Create(book models.Book) models.Book {
	book.ID = s.nextID
	s.nextID++
	s.books = append(s.books, book)
	return book
}

// Update replaces a book by ID
func (s *BookStore) Update(id int, book models.Book) error {
	for i := range s.books {
		if s.books[i].ID == id {
			book.ID = id // Preserve ID
			s.books[i] = book
			return nil
		}
	}
	return fmt.Errorf("book not found")
}

// Delete removes a book by ID
func (s *BookStore) Delete(id int) error {
	for i := range s.books {
		if s.books[i].ID == id {
			s.books = append(s.books[:i], s.books[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("book not found")
}
