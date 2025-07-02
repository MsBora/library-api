package storage

import (
	"database/sql"
	"fmt"
	"library-api/models"

	_ "github.com/jackc/pgx/v5/stdlib" //standard library database driver
)

// Storage Handles all database operations
type Storage struct {
	db *sql.DB
}

func NewStorage(connectionStr string) (*Storage, error) {
	db, err := sql.Open("pgx", connectionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open Database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}

	return &Storage{db: db}, nil
}

/*
	The *sql.DB object from Go's standard library is not a single database connection.
	It's a connection pool—a managed collection of database connections.
	When your application needs to talk to the database, it borrows a connection from the pool and returns it when done.
	This is extremely important for performance because creating a new database connection for every request is very slow.
	The connection pool handles all this for you automatically.
*/

func (s *Storage) CreateBook(book *models.Book) error {
	query := `INSERT INTO books (title, author, isbn, status) VALUES($1,$2,$3,$4) RETURNING id`

	err := s.db.QueryRow(query, &book.Title, &book.Author, &book.ISBN, &book.Status).Scan(&book.Id)
	if err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}

	return nil
}

func (s *Storage) GetBooks() ([]*models.Book, error) {
	query := `SELECT * FROM books`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get books: %w", err)
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		book := &models.Book{}
		if err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.ISBN, &book.Status); err != nil {
			return nil, fmt.Errorf("failed to scan book row: %w", err)
		}
		books = append(books, book)
	}
	return books, nil
}

func (s *Storage) GetBookbyId(id int) (*models.Book, error) {
	query := `SELECT * FROM books where id = $1`
	book := &models.Book{}
	err := s.db.QueryRow(query, id).Scan(&book.Id, &book.Title, &book.Author, &book.ISBN, &book.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Book not found")
		}
		return nil, fmt.Errorf("failed to get book by Id: %w", err)
	}

	return book, nil
}

func (s *Storage) UpdateBook(id int, book *models.Book) error {
	query := `Update books SET title = $1, author = $2, isbn = $3, status = $4 WHERE id = $5`

	_, err := s.db.Exec(query, &book.Title, &book.Author, &book.ISBN, &book.Status, id)
	if err != nil {
		return fmt.Errorf("Failed to update book: %w", err)
	}

	return nil
}

func (s *Storage) DeleteBook(id int) error {
	query := `DELETE FROM books WHERE id = $1`

	_, err := s.db.Exec(query, id)

	if err != nil {
		fmt.Errorf("Failed to delete the book: %w", err)
	}

	return nil
}

func (s *Storage) Close() {
	s.db.Close()
}
