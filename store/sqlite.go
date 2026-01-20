package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/favxlaw/models"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore manages books in SQLite database
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLite store
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	// Open database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create table if it doesn't exist
	err = RunMigrations(db)
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

// createTable creates the books table if it doesn't exist
func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'to_read',
		category TEXT,
		notes TEXT,
		start_date DATETIME NOT NULL,
		end_date DATETIME
	);`

	_, err := db.Exec(query)
	return err
}

// GetAll returns all books
func (s *SQLiteStore) GetAll() []models.Book {
	query := `
		SELECT id, title, author, status, category, notes, start_date, end_date 
		FROM books
		ORDER BY id DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return []models.Book{} // Return empty slice on error
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		book, err := scanBook(rows)
		if err != nil {
			continue // Skip invalid rows
		}
		books = append(books, book)
	}

	return books
}

// GetByID finds a book by ID
func (s *SQLiteStore) GetByID(id int) (*models.Book, error) {
	query := `
		SELECT id, title, author, status, category, notes, start_date, end_date 
		FROM books 
		WHERE id = ?
	`

	row := s.db.QueryRow(query, id)
	book, err := scanBookRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("book not found")
		}
		return nil, err
	}

	return &book, nil
}

// Create adds a new book and returns it with the generated ID
func (s *SQLiteStore) Create(book models.Book) models.Book {
	query := `
		INSERT INTO books (title, author, status, category, notes, start_date, end_date)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	// Convert end_date to proper format
	var endDate interface{}
	if book.EndDate != nil {
		endDate = book.EndDate.Format(time.RFC3339)
	}

	result, err := s.db.Exec(
		query,
		book.Title,
		book.Author,
		book.Status,
		book.Category,
		book.Notes,
		book.StartDate.Format(time.RFC3339),
		endDate,
	)

	if err != nil {
		return book // Return original book on error
	}

	// Get the auto-generated ID
	id, err := result.LastInsertId()
	if err != nil {
		return book
	}

	book.ID = int(id)
	return book
}

// Update replaces a book by ID
func (s *SQLiteStore) Update(id int, book models.Book) error {
	query := `
		UPDATE books 
		SET title = ?, author = ?, status = ?, category = ?, notes = ?, start_date = ?, end_date = ?
		WHERE id = ?
	`

	// Convert end_date to proper format
	var endDate interface{}
	if book.EndDate != nil {
		endDate = book.EndDate.Format(time.RFC3339)
	}

	result, err := s.db.Exec(
		query,
		book.Title,
		book.Author,
		book.Status,
		book.Category,
		book.Notes,
		book.StartDate.Format(time.RFC3339),
		endDate,
		id,
	)

	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

// Delete removes a book by ID
func (s *SQLiteStore) Delete(id int) error {
	query := `DELETE FROM books WHERE id = ?`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// Helper functions

// scanBook scans a row from Rows into a Book struct
func scanBook(rows *sql.Rows) (models.Book, error) {
	var book models.Book
	var startDateStr string
	var endDateStr sql.NullString // sql.NullString handles NULL values

	err := rows.Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Status,
		&book.Category,
		&book.Notes,
		&startDateStr,
		&endDateStr,
	)

	if err != nil {
		return book, err
	}

	// Parse start_date
	book.StartDate, _ = time.Parse(time.RFC3339, startDateStr)

	// Parse end_date if not NULL
	if endDateStr.Valid {
		endDate, _ := time.Parse(time.RFC3339, endDateStr.String)
		book.EndDate = &endDate
	}

	return book, nil
}

// scanBookRow scans a single row from QueryRow
func scanBookRow(row *sql.Row) (models.Book, error) {
	var book models.Book
	var startDateStr string
	var endDateStr sql.NullString

	err := row.Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Status,
		&book.Category,
		&book.Notes,
		&startDateStr,
		&endDateStr,
	)

	if err != nil {
		return book, err
	}

	// Parse dates
	book.StartDate, _ = time.Parse(time.RFC3339, startDateStr)

	if endDateStr.Valid {
		endDate, _ := time.Parse(time.RFC3339, endDateStr.String)
		book.EndDate = &endDate
	}

	return book, nil
}

// GetByFilters returns books matching the provided filters
func (s *SQLiteStore) GetByFilters(status, category, sortBy string) []models.Book {
	query := `
		SELECT id, title, author, status, category, notes, start_date, end_date 
		FROM books
		WHERE 1=1
	`
	args := []interface{}{}

	// Add filters
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}

	if category != "" {
		query += ` AND category = ?`
		args = append(args, category)
	}

	// Add sorting
	switch sortBy {
	case "title":
		query += ` ORDER BY title ASC`
	case "author":
		query += ` ORDER BY author ASC`
	case "date":
		query += ` ORDER BY start_date DESC`
	default:
		query += ` ORDER BY id DESC`
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return []models.Book{}
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		book, err := scanBook(rows)
		if err != nil {
			continue
		}
		books = append(books, book)
	}

	return books
}
