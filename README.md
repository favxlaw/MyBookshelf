# 📚 Personal Book Tracker

A book tracking REST API built in Go as a 3-month learning project, focusing on clean architecture, database fundamentals, and engineering best practices.

## 🚀 Quick Start

```bash
# Start the server
go run main.go

# The API will be available at http://localhost:8006
# Data persists in SQLite database: booktracker.db
```

## 📁 Project Structure

```
MyBookshelf/
├── main.go              # Application entry point & server setup
├── booktracker.db       # SQLite database (auto-created)
├── models/              # Data structures (Book, BookStatus)
│   └── book.go
├── library/             # Business logic (Phase 1 CRUD functions)
│   └── operations.go
├── handlers/            # HTTP request handlers
│   └── books.go
├── store/               # Data persistence layer
│   ├── sqlite.go        # SQLite implementation
│   └── migrations.go    # Database schema versioning
└── go.mod               # Go module dependencies
```

## 🔗 API Endpoints

### Books Collection
```bash
# List all books
GET /books

# Filter by status
GET /books?status=reading

# Filter by category
GET /books?category=Software%20Engineering

# Sort by title
GET /books?sort=title

# Combine filters
GET /books?status=reading&sort=title

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

# Filter by status
curl "http://localhost:8006/books?status=reading"

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

# Sort by title
curl "http://localhost:8006/books?sort=title"

# Sort by date (most recent first)
curl "http://localhost:8006/books?sort=date"
```

## 🏗️ Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────┐
│         main.go (The Manager)           │
│  - Creates SQLite connection            │
│  - Runs database migrations             │
│  - Creates bookHandler                  │
│  - Registers routes                     │
│  - Starts server                        │
└─────────────────┬───────────────────────┘
                  │
                  ↓
┌─────────────────────────────────────────┐
│    handlers/books.go (HTTP Layer)       │
│  - ServeHTTP() receives all requests    │
│  - getAllBooks()      ← GET /books      │
│  - createBook()       ← POST /books     │
│  - getBookByID()      ← GET /books/5    │
│  - updateBook()       ← PUT /books/5    │
│  - deleteBook()       ← DELETE /books/5 │
│  - Query param parsing & validation     │
└─────────────────┬───────────────────────┘
                  │
                  ↓
┌─────────────────────────────────────────┐
│    store/sqlite.go (Data Layer)         │
│  - GetAll()                             │
│  - GetByID()                            │
│  - GetByFilters() ← Filtering & sorting │
│  - Create()                             │
│  - Update()                             │
│  - Delete()                             │
│  - SQL query building                   │
└─────────────────┬───────────────────────┘
                  │
                  ↓
┌─────────────────────────────────────────┐
│       SQLite Database (Persistence)     │
│  - books table                          │
│  - schema_migrations table              │
│  - Auto-incrementing IDs                │
│  - Data survives server restarts        │
└─────────────────────────────────────────┘
```

### Design Principles
- **Separation of Concerns**: Each layer has a single responsibility
- **Dependency Injection**: Components receive dependencies via interfaces
- **Interface-Based Design**: Store interface allows swapping implementations
- **Database Abstraction**: Business logic doesn't know about SQL
- **Migration System**: Version-controlled schema changes

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

### Phase 3: Database Integration ✅ (Complete)
**What I Built:**
- SQLite database integration with `database/sql`
- Migration system for schema version control
- SQL-based CRUD operations with prepared statements
- Advanced filtering and sorting at database level
- Connection management and resource cleanup
- Data persistence across server restarts

**What I Learned:**
- SQL fundamentals (CREATE, INSERT, SELECT, UPDATE, DELETE)
- Database drivers in Go (`database/sql` package)
- Prepared statements and SQL injection prevention
- Handling NULL values with `sql.NullString`
- Auto-increment primary keys
- Database migrations concept and implementation
- Query parameter parsing and dynamic SQL
- Resource management with `defer`
- Error handling for database operations

---

### Phase 4: Production Ready (Coming Soon)
- Environment configuration (dev, staging, prod)
- Structured logging (request tracking, error logs)
- Middleware (CORS, authentication, rate limiting)
- Unit and integration tests
- Docker containerization
- Deployment basics

---

##  What I've Learned So Far

### Go Fundamentals
- Package system and imports
- Structs and custom types
- Pointers vs values
- Slices and arrays
- Error handling with `error` interface
- Methods and receivers
- Interfaces and polymorphism
- `defer` for resource cleanup

### HTTP & APIs
- Request/response flow
- REST conventions
- JSON marshaling/unmarshaling
- HTTP status codes
- URL routing and parameters
- Query parameters

### Database Fundamentals
- SQL syntax and operations
- Database drivers and connections
- Prepared statements
- Transaction concepts
- Schema migrations
- NULL handling
- Auto-increment IDs
- Resource management

### Software Engineering
- CRUD operations
- Input validation
- Clean architecture
- Separation of concerns
- Dependency injection
- Interface-based design
- Error propagation
- Defensive programming

### Development Practices
- Incremental learning (one concept at a time)
- Breaking complex problems into steps
- Thinking in abstractions
- Writing maintainable code
- Version controlling schema changes



**Status**: Phase 3 Complete ✅  
**Next Milestone**: Production-Ready Features  
**Current Focus**: Understanding database fundamentals and clean architecture patterns
