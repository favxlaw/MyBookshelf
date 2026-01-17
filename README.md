# 📚 Personal Book Tracker

A book tracking REST API built in Go as a learning project, focusing on clean architecture and engineering fundamentals.

## 🚀 Quick Start

```bash
# Start the server
go run main.go

# The API will be available at http://localhost:8006
```

## 📁 Project Structure

```
MyBookshelf/
├── main.go              # Application entry point & server setup
├── models/              # Data structures (Book, BookStatus)
│   └── book.go
├── library/             # Business logic (Phase 1 CRUD functions)
│   └── operations.go
├── handlers/            # HTTP request handlers
│   └── books.go
├── store/               # Data persistence layer
│   └── memory.go
└── go.mod               # Go module dependencies
```

## 🔗 API Endpoints

### Books Collection
```bash
# List all books
GET /books

# Add a new book
POST /books
Content-Type: application/json
{
  "Title": "Atomic Habits",
  "Author": "James Clear",
  "Status": "reading",
  "Category": "Self-Help"
}
```

### Single Book Operations
```bash
# Get specific book
GET /books/{id}

# Update a book
PUT /books/{id}
Content-Type: application/json
{
  "Title": "Clean Code",
  "Author": "Robert C. Martin",
  "Status": "finished",
  "Category": "Software Engineering"
}

# Delete a book
DELETE /books/{id}
```

## 📖 Usage Examples

```bash
# List all books
curl http://localhost:8006/books

# Add a new book
curl -X POST http://localhost:8006/books \
  -H "Content-Type: application/json" \
  -d '{
    "Title": "The Pragmatic Programmer",
    "Author": "Hunt & Thomas",
    "Status": "reading",
    "Category": "Software Engineering"
  }'

# Get book with ID 1
curl http://localhost:8006/books/1

# Update book status
curl -X PUT http://localhost:8006/books/2 \
  -H "Content-Type: application/json" \
  -d '{
    "Title": "Clean Code",
    "Author": "Robert C. Martin",
    "Status": "finished",
    "Category": "Software Engineering"
  }'

# Delete a book
curl -X DELETE http://localhost:8006/books/3
```

## 🏗️ Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────┐
│         main.go (The Manager)           │
│  - Creates bookStore                    │
│  - Creates bookHandler                  │
│  - Registers routes                     │
│  - Starts server                        │
└─────────────────┬───────────────────────┘
                  │
                  ↓
┌─────────────────────────────────────────┐
│    handlers/books.go (The Router)       │
│  - ServeHTTP() receives all requests    │
│  - getAllBooks()      ← GET /books      │
│  - createBook()       ← POST /books     │
│  - getBookByID()      ← GET /books/5    │
│  - updateBook()       ← PUT /books/5    │
│  - deleteBook()       ← DELETE /books/5 │
└─────────────────┬───────────────────────┘
                  │
                  ↓
┌─────────────────────────────────────────┐
│    store/memory.go (The Database)       │
│  - GetAll()                             │
│  - GetByID()                            │
│  - Create()                             │
│  - Update()                             │
│  - Delete()                             │
└─────────────────────────────────────────┘
```

### Design Principles
- **Separation of Concerns**: Each package has a single, well-defined responsibility
- **Dependency Injection**: Components receive their dependencies rather than creating them
- **Interface-Based Design**: `http.Handler` interface for clean HTTP handling
- **Encapsulation**: Data access through store methods, not global variables

## ✅ Progress

### Phase 1: Core Data Operations ✅ (Complete)
**What I Built:**
- Book struct with proper types and validation
- CRUD operations (Create, Read, Update, Delete)
- Filter functions by status and category
- In-memory storage using slices
- Package organization (models, library)

**What I Learned:**
- Data modeling with structs
- Type safety with custom types and constants
- Slices and collections
- Pointers for data mutation
- Functions and package organization

---

### Phase 2: REST API ✅ (Complete)
**What I Built:**
- Full REST API with proper HTTP methods
- JSON request/response handling
- URL parameter extraction
- Input validation with helpful error messages
- Clean architecture with separated concerns
- Professional project structure

**What I Learned:**
- HTTP fundamentals (request/response cycle)
- REST principles and conventions
- JSON encoding/decoding in Go
- HTTP methods (GET, POST, PUT, DELETE)
- Status codes (200, 201, 204, 400, 404, 405)
- Structs with methods (receivers)
- The `http.Handler` interface
- Dependency injection pattern
- Manual routing with `net/http`

**Key Concepts Applied:**
- Request validation and error handling
- RESTful URL design (`/books` vs `/books/{id}`)
- Separation of HTTP logic from business logic
- In-memory data store with clean interface

---

### Phase 3: Database Integration (Coming Soon)
- Replace in-memory storage with PostgreSQL/SQLite
- SQL fundamentals (SELECT, INSERT, UPDATE, DELETE)
- Database migrations
- Connection pooling
- Data persistence across server restarts


## 📝 Notes

- Currently using in-memory storage (data resets on restart)
- Using Go's standard `net/http` library (no frameworks)
- Manual routing to understand HTTP fundamentals
- Will add database in Phase 3 for persistence

