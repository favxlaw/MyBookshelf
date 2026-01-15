package models

import "time"

// Book represents a book in the library
type Book struct {
	ID        int
	Title     string
	Author    string
	Status    BookStatus
	Category  string
	Notes     string
	StartDate time.Time
	EndDate   *time.Time
}

// BookStatus represents the reading status of a book
type BookStatus string

const (
	StatusToRead    BookStatus = "to_read"
	StatusReading   BookStatus = "reading"
	StatusFinished  BookStatus = "finished"
	StatusAbandoned BookStatus = "abandoned"
)
