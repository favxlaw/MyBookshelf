package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/favxlaw/models"
	"github.com/favxlaw/store"
)

// BookHandler handles all book-related HTTP requests
type BookHandler struct {
	store *store.BookStore
}

// NewBookHandler creates a new book handler
func NewBookHandler(store *store.BookStore) *BookHandler {
	return &BookHandler{
		store: store,
	}
}

// ServeHTTP implements http.Handler interface
func (h *BookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if path has an ID
	if r.URL.Path != "/books" && r.URL.Path != "/books/" {
		h.handleSingleBook(w, r)
		return
	}

	// Collection operations
	switch r.Method {
	case http.MethodGet:
		h.getAllBooks(w, r)
	case http.MethodPost:
		h.createBook(w, r)
	default:
		errorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleSingleBook routes single book operations
func (h *BookHandler) handleSingleBook(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getBookByID(w, r, id)
	case http.MethodPut:
		h.updateBook(w, r, id)
	case http.MethodDelete:
		h.deleteBook(w, r, id)
	default:
		errorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getAllBooks handles GET /books
func (h *BookHandler) getAllBooks(w http.ResponseWriter, _ *http.Request) {
	books := h.store.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// getBookByID handles GET /books/{id}
func (h *BookHandler) getBookByID(w http.ResponseWriter, _ *http.Request, id int) {
	book, err := h.store.GetByID(id)
	if err != nil {
		errorResponse(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// createBook handles POST /books
func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	var newBook models.Book

	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		errorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	err = validateBook(newBook)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set defaults
	newBook.StartDate = time.Now()
	if newBook.Status == "" {
		newBook.Status = models.StatusToRead
	}

	// Create in store
	created := h.store.Create(newBook)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// updateBook handles PUT /books/{id}
func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request, id int) {
	// Get existing book
	existingBook, err := h.store.GetByID(id)
	if err != nil {
		errorResponse(w, "Book not found", http.StatusNotFound)
		return
	}

	// Decode update
	var updatedBook models.Book
	err = json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		errorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate
	err = validateBook(updatedBook)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Preserve certain fields
	updatedBook.StartDate = existingBook.StartDate

	// Handle EndDate logic
	if (updatedBook.Status == models.StatusFinished || updatedBook.Status == models.StatusAbandoned) &&
		existingBook.EndDate == nil {
		now := time.Now()
		updatedBook.EndDate = &now
	} else if updatedBook.Status == models.StatusReading || updatedBook.Status == models.StatusToRead {
		updatedBook.EndDate = nil
	} else {
		updatedBook.EndDate = existingBook.EndDate
	}

	// Update in store
	err = h.store.Update(id, updatedBook)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBook)
}

// deleteBook handles DELETE /books/{id}
func (h *BookHandler) deleteBook(w http.ResponseWriter, _ *http.Request, id int) {
	err := h.store.Delete(id)
	if err != nil {
		errorResponse(w, "Book not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper functions

func extractID(path string) (int, error) {
	idStr := strings.TrimPrefix(path, "/books/")
	if idStr == "" || idStr == path {
		return 0, fmt.Errorf("no ID provided")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format")
	}

	return id, nil
}

func validateBook(book models.Book) error {
	if book.Title == "" {
		return fmt.Errorf("title is required")
	}

	if book.Author == "" {
		return fmt.Errorf("author is required")
	}

	if book.Status != "" {
		validStatuses := []models.BookStatus{
			models.StatusToRead,
			models.StatusReading,
			models.StatusFinished,
			models.StatusAbandoned,
		}

		isValid := false
		for _, validStatus := range validStatuses {
			if book.Status == validStatus {
				isValid = true
				break
			}
		}

		if !isValid {
			return fmt.Errorf("status must be one of: to_read, reading, finished, abandoned")
		}
	}

	return nil
}

func errorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
