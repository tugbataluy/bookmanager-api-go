package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "postgres"
	dbname = "bookmanager"
)

var DB *sql.DB

func InitDB() (*sql.DB, error){
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}


	fmt.Println("Successfully connected to PostgreSQL database")
	createTables()
	return DB, nil
	
}


func createTables() {
	createBooksTable := `
	CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		published_date DATE NOT NULL,
		edition INTEGER NOT NULL DEFAULT 1,
		description TEXT,
		genre VARCHAR(100),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	
	createCollectionsTable := `
	CREATE TABLE IF NOT EXISTS collections (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	createCollectionBooksTable := `
	CREATE TABLE IF NOT EXISTS collection_books (
		collection_id INTEGER REFERENCES collections(id) ON DELETE CASCADE,
		book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (collection_id, book_id)
	);`
	
	_, err := DB.Exec(createBooksTable)
	if err != nil {
		log.Fatal("Couldn't create books table:", err)
	}

	_, err = DB.Exec(createCollectionsTable)
	if err != nil {
		log.Fatal("Couldn't create collections table:", err)
	}

	_, err = DB.Exec(createCollectionBooksTable)
	if err != nil {
		log.Fatal("Couldn't create collection_books table:", err)
	}

	createIndexes()

}

func createIndexes() {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);",
		"CREATE INDEX IF NOT EXISTS idx_books_genre ON books(genre);",
		"CREATE INDEX IF NOT EXISTS idx_books_published_date ON books(published_date);",
	}

	for _, index := range indexes {
		_, err := DB.Exec(index)
		if err != nil {
			log.Printf("Couldn't create index: %v", err)
		}
	}
}
