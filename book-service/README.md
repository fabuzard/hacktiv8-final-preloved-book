# Book Service

Book management microservice for the used book marketplace.

## Features

- Create, read, update, delete books
- Filter books by category
- Seller-specific book management
- JWT authentication ready (placeholder middleware)

## API Endpoints

### Books
- `POST /api/v1/books` - Create a new book
- `GET /api/v1/books` - Get all books (with optional category filter)
- `GET /api/v1/books/my` - Get books for authenticated seller
- `GET /api/v1/books/:id` - Get book by ID
- `PUT /api/v1/books/:id` - Update book (seller only)
- `DELETE /api/v1/books/:id` - Delete book (seller only)

## Setup

1. Copy environment variables:
   ```bash
   cp .env.example .env
   ```

2. Update database configuration in `.env`

3. Run the service:
   ```bash
   go run main.go
   ```

## Database Migration

The service automatically migrates the Book model on startup using GORM AutoMigrate.