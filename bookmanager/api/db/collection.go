package db

import (
	"bookmanager/api/models"
	"database/sql"
	"fmt"
)

type CollectionDB struct {
	DB *sql.DB
}

func NewCollection(db *sql.DB) *CollectionDB {
	return &CollectionDB{DB: db}
}

func (c* CollectionDB) CreateCollection(collection *models.CollectionRequest) (*models.Collection, error) {
	var newCollection models.Collection
	query := `
	INSERT INTO collections (name, description)
	VALUES ($1, $2)
	RETURNING id, name, description, created_at, updated_at`

	err := DB.QueryRow(
		query,
		collection.Name,
		collection.Description,
	).Scan(
		&newCollection.ID,
		&newCollection.Name,
		&newCollection.Description,
		&newCollection.CreatedAt,
		&newCollection.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create collection: %v", err)
	}

	return  &newCollection, nil
}

func (c* CollectionDB) GetCollection(id int) (*models.Collection, error) {
	var collection models.Collection
	query := `
	SELECT id, name, description, created_at, updated_at
	FROM collections
	WHERE id = $1`
	
	err := DB.QueryRow(query, id).Scan(
		&collection.ID,
		&collection.Name,
		&collection.Description,
		&collection.CreatedAt,
		&collection.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("collection not found")
		}
		return  nil, fmt.Errorf("failed to get collection :%v", err)
	}

	return &collection, nil
}


func (c* CollectionDB) UpdateCollection(id int, collection *models.CollectionRequest) (*models.Collection, error) {
	var updatedCollection models.Collection
	query := `
	UPDATE collections
	SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
	WHERE id = $3
	RETURNING id, name, description, created_at, updated_at`

	err := DB.QueryRow(
		query,
		collection.Name,
		collection.Description,
		id,
	).Scan(
		&updatedCollection.ID,
		&updatedCollection.Name,
		&updatedCollection.Description,
		&updatedCollection.CreatedAt,
		&updatedCollection.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("collection not found")
		}
		return nil, fmt.Errorf("failed to update collection: %v", err)
	}
	return  &updatedCollection, nil
}

func (c *CollectionDB) PatchCollection(id int, patch *models.CollectionRequest) (*models.Collection, error) {
    current, err := c.GetCollection(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get current collection: %w", err)
    }

    merged := mergeCollectionWithPatch(current, patch)
    return c.UpdateCollection(id, merged)
}

func mergeCollectionWithPatch(current *models.Collection, patch *models.CollectionRequest) *models.CollectionRequest {
    merged := &models.CollectionRequest{
        Name:        current.Name,
        Description: current.Description,
    }

    if patch.Name != "" {
        merged.Name = patch.Name
    }
    if patch.Description != "" {
        merged.Description = patch.Description
    }

    return merged
}

func (c* CollectionDB) DeleteCollection(id int) error {
	query := `DELETE FROM collections WHERE id = $1`
	result, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("collection not found")
	}

	return nil
}

func (c *CollectionDB) ListCollections(where, groupBy, orderBy, limit, offset string) (interface{}, error) {
    baseQuery := `
        SELECT id, name, description, created_at, updated_at
        FROM collections`
    
    if groupBy != "" {
        query := fmt.Sprintf(`
            SELECT %s as group_key, COUNT(*) as count
            FROM collections
            GROUP BY %s`, groupBy, groupBy)
        
        rows, err := DB.Query(query)
        if err != nil {
            return nil, fmt.Errorf("failed to list grouped collections: %v", err)
        }
        defer rows.Close()

        groups := make(map[string]int)
        for rows.Next() {
            var key string
            var count int
            if err := rows.Scan(&key, &count); err != nil {
                return nil, fmt.Errorf("failed to scan group: %v", err)
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
        baseQuery += " ORDER BY name"
    }
    if limit != "" {
        baseQuery += " LIMIT " + limit
    }
    if offset != "" {
        baseQuery += " OFFSET " + offset
    }

    rows, err := DB.Query(baseQuery)
    if err != nil {
        return nil, fmt.Errorf("failed to list collections: %v", err)
    }
    defer rows.Close()

    var collections []models.Collection
    for rows.Next() {
        var collection models.Collection
        err := rows.Scan(
            &collection.ID,
            &collection.Name,
            &collection.Description,
            &collection.CreatedAt,
            &collection.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan collection: %v", err)
        }
        collections = append(collections, collection)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error after scanning collections: %v", err)
    }

    return collections, nil
}

func (c* CollectionDB) AddBookToCollection(collectionID, bookID int) error {
	query := `
	INSERT INTO collection_books (collection_id, book_id)
	VALUES ($1, $2)
	ON CONFLICT (collection_id, book_id) DO NOTHING`

	result, err := DB.Exec(query, collectionID, bookID)
	if err != nil {
		return fmt.Errorf("failed to add book to collection: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return  fmt.Errorf("failed to check rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book already exists in collection")
	}

	return nil
}

func (c* CollectionDB) RemoveBookFromCollection(collectionID, bookID int) error {
	query :=  `
	DELETE FROM collection_books
	WHERE collection_id = $1 AND book_id = $2`

	result, err := DB.Exec(query, collectionID, bookID)
	if err != nil {
		return fmt.Errorf("failed to remove book from collections: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return  fmt.Errorf("failed to check rows affected: %v", err)
	}

	
	if rowsAffected == 0 {
		return fmt.Errorf("book not found in collection")
	}

	return nil

}


func (c* CollectionDB) ListBooksInCollection(collectionID int) ([]models.Book, error) {
	query := `
	SELECT b.id, b.title, b.author, b.published_date, b.edition, b.description, b.genre, b.created_at, b.updated_at
	FROM books b
	JOIN collection_books cb ON b.id = cb.book_id
	WHERE cb.collection_id = $1
	ORDER BY b.title`

	rows, err := DB.Query(query, collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list books in collection: %v", err)
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
