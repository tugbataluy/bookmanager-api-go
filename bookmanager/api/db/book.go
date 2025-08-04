package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"bookmanager/api/models"
)

type BookDB struct {
	DB *sql.DB
}

type GroupedResult struct {
	GroupKey string
	Books    []models.Book
}

func NewBook(db *sql.DB) *BookDB {
	return &BookDB{DB: db}
}

func (b *BookDB) CreateBook(book *models.BookRequest) (*models.Book, error) {
	publishedDate, err := time.Parse("2006-01-02", book.PublishedDate)
	if err != nil {
		return nil, fmt.Errorf("invalid published date format: %v", err)
	}

	var newBook models.Book
	query := `
		INSERT INTO books (title, author, published_date, edition, description, genre)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, author, published_date, edition, description, genre, created_at, updated_at`

	err = DB.QueryRow(
		query,
		book.Title,
		book.Author,
		publishedDate,
		book.Edition,
		book.Description,
		book.Genre,
	).Scan(
		&newBook.ID,
		&newBook.Title,
		&newBook.Author,
		&newBook.PublishedDate,
		&newBook.Edition,
		&newBook.Description,
		&newBook.Genre,
		&newBook.CreatedAt,
		&newBook.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create book %v", err)
	}

	return &newBook, nil
}

func (b *BookDB) GetBook(id int) (*models.Book, error) {
	var book models.Book
	query := `
	SELECT id, title, author, published_date, edition, description, genre, created_at, updated_at
	FROM books
	WHERE id = $1
	`

	err := DB.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.PublishedDate,
		&book.Edition,
		&book.Description,
		&book.Genre,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("book not found")
		}
		return nil, fmt.Errorf("failed to get book: %v", err)
	}

	return &book, nil
}

func (b *BookDB) UpdateBook(id int, book *models.BookRequest) (*models.Book, error) {
	publishedDate, err := time.Parse("2006-01-02", book.PublishedDate)
	if err != nil {
		return nil, fmt.Errorf("invalid published date format: %v", err)
	}

	var updatedBook models.Book
	query := `
	UPDATE books
	SET title = $1, author = $2, published_date = $3, edition = $4, 
	    description = $5, genre = $6, updated_at = CURRENT_TIMESTAMP
	WHERE id = $7
	RETURNING id, title, author, published_date, edition, description, genre, created_at, updated_at`
	err = DB.QueryRow(
		query,
		book.Title,
		book.Author,
		publishedDate,
		book.Edition,
		book.Description,
		book.Genre,
		id,
	).Scan(
		&updatedBook.ID,
		&updatedBook.Title,
		&updatedBook.Author,
		&updatedBook.PublishedDate,
		&updatedBook.Edition,
		&updatedBook.Description,
		&updatedBook.Genre,
		&updatedBook.CreatedAt,
		&updatedBook.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("book not found")
		}
		return nil, fmt.Errorf("failed to update book: %v", err)
	}

	return &updatedBook, nil
}

func (b *BookDB) PatchBook(id int, patch *models.BookRequest) (*models.Book, error) {
	current, err := b.GetBook(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get current book: %w", err)
	}

	merged := mergeBookWithPatch(current, patch)

	return b.UpdateBook(id, merged)
}

func mergeBookWithPatch(current *models.Book, patch *models.BookRequest) *models.BookRequest {
	currentDateParts := strings.Split(current.PublishedDate, "T")
	currentDateOnly := currentDateParts[0]
	merged := &models.BookRequest{
		Title:         current.Title,
		Author:        current.Author,
		PublishedDate: currentDateOnly,
		Edition:       current.Edition,
		Description:   current.Description,
		Genre:         current.Genre,
	}

	if patch.Title != "" {
		merged.Title = patch.Title
	}
	if patch.Author != "" {
		merged.Author = patch.Author
	}
	if patch.PublishedDate != "" {
		merged.PublishedDate = patch.PublishedDate
	}
	if patch.Edition != 0 {
		merged.Edition = patch.Edition
	}
	if patch.Description != "" {
		merged.Description = patch.Description
	}
	if patch.Genre != "" {
		merged.Genre = patch.Genre
	}

	return merged
}

func (b *BookDB) DeleteBook(id int) error {
	query := `DELETE FROM books WHERE id = $1`
	result, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book : %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

func (b *BookDB) ListBooks(where, groupBy, orderBy, limit, offset string) (interface{}, error) {
	baseQuery := `
        SELECT id, title, author, published_date, edition, 
               description, genre, created_at, updated_at 
        FROM books`

	if groupBy != "" {
		query := fmt.Sprintf(`
            SELECT %s as group_key, COUNT(*) as count
            FROM books
            GROUP BY %s`, groupBy, groupBy)

		rows, err := DB.Query(query)
		if err != nil {
			return nil, fmt.Errorf("failed to list grouped books: %v", err)
		}
		defer rows.Close()

		groups := make(map[string]int)
		for rows.Next() {
			var key string
			var count int
			if err := rows.Scan(&key, &count); err != nil {
				return nil, fmt.Errorf("failed to scan book group: %v", err)
			}
			groups[key] = count
		}
		return groups, nil
	}

	if where != "" {
		baseQuery += " WHERE " + where
	}

	if orderBy != "" {
		baseQuery += " ORDER BY " + orderBy
	} else {
		baseQuery += " ORDER BY title" // Default ordering by title
	}

	if limit != "" {
		baseQuery += " LIMIT " + limit
	}

	if offset != "" {
		baseQuery += " OFFSET " + offset
	}

	rows, err := DB.Query(baseQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to list books: %v", err)
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.PublishedDate,
			&book.Edition,
			&book.Description,
			&book.Genre,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan book: %v", err)
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning books: %v", err)
	}

	return books, nil
}
