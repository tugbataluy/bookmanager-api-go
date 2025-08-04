package main

import (
	"bookmanager/api/db"
	"bookmanager/api/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dbConn, err:= db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbConn.Close()

	bookHandler := handlers.NewBookHandler(&db.BookDB{})
	collectionHandler := handlers.NewCollectionHandler(&db.CollectionDB{})

	http.HandleFunc("/api/v1/books", bookHandler.HandleBooks)
	http.HandleFunc("/api/v1/books/", bookHandler.HandleBook)
	http.HandleFunc("/api/v1/collections", collectionHandler.HandleCollections)
	http.HandleFunc("/api/v1/collections/{id}", collectionHandler.HandleCollection)
	http.HandleFunc("/api/v1/collections-books/", collectionHandler.HandleCollectionBooksRoutes)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server {
		Addr: fmt.Sprintf(":%s", port),
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			log.Fatalf("Server error: %v", err)
		}
	}()
	<-done
	log.Println("Server stopped")
}
