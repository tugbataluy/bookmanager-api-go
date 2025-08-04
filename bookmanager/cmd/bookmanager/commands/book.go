package commands

import (
	"bookmanager/api/models"
	"bookmanager/cmd/bookmanager/api"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func HandleBookCommand(client *api.APIClient, args []string) {
	if len(args) < 1 {
		printBookHelp()
		os.Exit(1)
	}

	switch args[0] {
	case "create":
		createBook(client, args[1:])
	case "list":
		listBooks(client, args[1:])
	case "get":
		getBook(client, args[1:])
	case "update":
		updateBook(client, args[1:])
	case "patch":
		patchBook(client, args[1:])
	case "delete":
		deleteBook(client, args[1:])
	case "help":
		printBookHelp()
	default:
		fmt.Printf("Unknown book command: %s\n", args[0])
		printBookHelp()
		os.Exit(1)
	}
}

func printBookHelp() {
	fmt.Println(`Usage: bookmanager book <command> [options]

Commands:
  create      Add a new book
  list        List books with optional filters
  get         Get details of a specific book
  update      Update a book's information (all fields required)
  patch	      Update a books' information partially (only updates provided fields)
  delete      Remove a book from the system
  help        Show this help message

List Options:
  --author          Filter by author
  --genre           Filter by genre
  --published-after Filter by publication date (after)
  --published-before Filter by publication date (before)
  --where           SQL-like WHERE clause (e.g., "title LIKE '%Hobbit%' AND edition > 1")
  --order-by        SQL-like ORDER BY clause (e.g., "published_date DESC")
  --group-by        SQL-like GROUP BY clause (e.g., "genre")
  --limit           Limit number of results
  --offset          Offset for pagination

Examples:
  bookmanager book create --title "The Hobbit" --author "J.R.R. Tolkien" --published-date "1937-09-21"
  bookmanager book list --where "genre = 'Fantasy' AND published_date > '1950-01-01'"
  bookmanager book list --group-by "author"`)
}

func createBook(client *api.APIClient, args []string) {
	fs := flag.NewFlagSet("book create", flag.ExitOnError)
	title := fs.String("title", "", "Book title (required)")
	author := fs.String("author", "", "Book author (required)")
	publishedDate := fs.String("published-date", "", "Publication Date (YYYY-MM-DD) (required)")
	edition := fs.Int("edition", 1, "Edition number")
	description := fs.String("description", "", "Book description")
	genre := fs.String("genre", "", "Book genre")
	fs.Parse(args)

	if *title == "" || *author == "" || *publishedDate == "" {
		fmt.Println("Title, author, and published-date are required")
		fs.PrintDefaults()
		os.Exit(1)
	}

	if _, err := time.Parse("2006-01-02", *publishedDate); err != nil {
		fmt.Println("Invalid date format. Please use YYYY-MM-DD")
		os.Exit(1)
	}

	book := map[string]interface{}{
		"title":          *title,
		"author":         *author,
		"published_date": *publishedDate,
		"edition":        *edition,
		"description":    *description,
		"genre":          *genre,
	}

	body, err := client.Post("/books", book)
	if err != nil {
		log.Fatalf("Error creating book %v", err)
	}

	var createdBook models.Book
	if err := json.Unmarshal(body, &createdBook); err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}
	fmt.Printf("Created book %d: %s by %s\n", createdBook.ID, createdBook.Title, createdBook.Author)
}

func getBook(client *api.APIClient, args []string) {
	if len(args) < 1 {
		fmt.Println("Book ID is required")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid book ID")
		os.Exit(1)
	}

	body, err := client.Get(fmt.Sprintf("/v1/books/%d", id), nil)
	if err != nil {
		log.Fatalf("Error getting book: %v", err)
	}

	var book models.Book
	if err := json.Unmarshal(body, &book); err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	fmt.Printf("Book #%d\n", book.ID)
	fmt.Printf("Title: %s\n", book.Title)
	fmt.Printf("Author: %s\n", book.Author)
	fmt.Printf("Published Date: %s\n", book.PublishedDate)
	fmt.Printf("Edition: %d\n", book.Edition)
	if book.Description != "" {
		fmt.Printf("Description: %s\n", book.Description)
	}
	if book.Genre != "" {
		fmt.Printf("Genre: %s\n", book.Genre)
	}
}

func listBooks(client *api.APIClient, args []string) {
	fs := flag.NewFlagSet("book list", flag.ExitOnError)
	where := fs.String("where", "", "SQL-like WHERE clause")
	groupBy := fs.String("group-by", "", "SQL-like GROUP BY clause")
	orderBy := fs.String("order-by", "", "SQL-like ORDER BY clause")
	limit := fs.Int("limit", 0, "Limit number of results")
	offset := fs.Int("offset", 0, "Offset for pagination")

	author := fs.String("author", "", "Filter by author")
	genre := fs.String("genre", "", "Filter by genre")
	publishedAfter := fs.String("published-after", "", "Filter by publication date (after)")
	publishedBefore := fs.String("published-before", "", "Filter by publication date (before)")

	if err := fs.Parse(args); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	if *limit < 0 || *offset < 0 {
		log.Fatal("Limit and offset must be positive numbers")
	}

	if *publishedAfter != "" {
		if _, err := time.Parse("2006-01-02", *publishedAfter); err != nil {
			log.Fatalf("Invalid published-after date format: %v", err)
		}
	}
	if *publishedBefore != "" {
		if _, err := time.Parse("2006-01-02", *publishedBefore); err != nil {
			log.Fatalf("Invalid published-before date format: %v", err)
		}
	}

	params := client.BuildQueryParams(*where, *groupBy, *orderBy, *limit, *offset)

	if *author != "" {
		params["author"] = *author
	}
	if *genre != "" {
		params["genre"] = *genre
	}
	if *publishedAfter != "" {
		params["published_after"] = *publishedAfter
	}
	if *publishedBefore != "" {
		params["published_before"] = *publishedBefore
	}

	body, err := client.Get("/v1/books", params)
	if err != nil {
		log.Fatalf("API request failed: %v", err)
	}

	if *groupBy != "" {
		var result struct {
			Groups map[string]int `json:"groups"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatalf("Error parsing grouped response: %v", err)
		}

		if len(result.Groups) == 0 {
			fmt.Println("No grouped results found")
			return
		}

		fmt.Printf("Books grouped by %s:\n", *groupBy)
		for group, count := range result.Groups {
			fmt.Printf("- %s: %d\n", group, count)
		}
	} else {
		var result struct {
			Books []models.Book `json:"books"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			log.Printf("Raw response: %s", string(body))
			log.Fatalf("Failed to parse response: %v", err)
		}

		if len(result.Books) == 0 {
			fmt.Println("No books found")
			return
		}

		for _, book := range result.Books {
			fmt.Printf("%d: %s by %s (%s)\n",
				book.ID, book.Title, book.Author, book.PublishedDate)
			if book.Edition > 1 {
				fmt.Printf("   Edition: %d\n", book.Edition)
			}
			if book.Genre != "" {
				fmt.Printf("   Genre: %s\n", book.Genre)
			}
			if book.Description != "" {
				fmt.Printf("   Description: %s\n", book.Description)
			}
			fmt.Println()
		}
	}
}

func updateBook(client *api.APIClient, args []string) {
	fs := flag.NewFlagSet("book update", flag.ExitOnError)
	title := fs.String("title", "", "Book title (required)")
	author := fs.String("author", "", "Book author (required)")
	publishedDate := fs.String("published-date", "", "Publication date (YYYY-MM-DD) (required)")
	edition := fs.Int("edition", 0, "Edition number (required)")
	description := fs.String("description", "", "Book description (required)")
	genre := fs.String("genre", "", "Book genre (required)")

	if len(args) < 1 {
		fmt.Println("Book ID is required")
		fs.PrintDefaults()
		os.Exit(1)
	}

	fs.Parse(args[1:])

	if *title == "" || *author == "" || *publishedDate == "" || *edition == 0 || *description == "" || *genre == "" {
		fmt.Println("Error: All fields are required for update. Use 'patch' for partial updates.")
		fmt.Println("Required flags:")

		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid book ID")
		os.Exit(1)
	}

	if _, err := time.Parse("2006-01-02", *publishedDate); err != nil {
		fmt.Println("Invalid date format. Please use YYYY-MM-DD")
		os.Exit(1)
	}

	updateData := map[string]interface{}{
		"title":          *title,
		"author":         *author,
		"published_date": *publishedDate,
		"edition":        *edition,
		"description":    *description,
		"genre":          *genre,
	}

	body, err := client.Put(fmt.Sprintf("/books/%d", id), updateData)
	if err != nil {
		log.Fatalf("Error updating book: %v", err)
	}

	var updatedBook models.Book
	if err := json.Unmarshal(body, &updatedBook); err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	fmt.Printf("Updated book #%d: %s by %s\n", updatedBook.ID, updatedBook.Title, updatedBook.Author)
}

func patchBook(client *api.APIClient, args []string) {
	fs := flag.NewFlagSet("book patch", flag.ExitOnError)
	title := fs.String("title", "", "Update book title (optional)")
	author := fs.String("author", "", "Update book author (optional)")
	publishedDate := fs.String("published-date", "", "Update publication date (YYYY-MM-DD, optional)")
	edition := fs.Int("edition", 1, "Update edition number (optional, 1 to skip)")
	description := fs.String("description", "", "Update book description (optional)")
	genre := fs.String("genre", "", "Update book genre (optional)")

	if len(args) == 0 {
		fs.PrintDefaults()
		os.Exit(1)
	}

	if len(args) < 1 {
		fmt.Println("Book ID is required")
		os.Exit(1)
	}

	fs.Parse(args[1:])
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid book ID")
		os.Exit(1)
	}

	patchData := make(map[string]interface{})

	if *title != "" {
		patchData["title"] = *title
	}
	if *author != "" {
		patchData["author"] = *author
	}
	if *publishedDate != "" {
		patchData["published_date"] = *publishedDate
	}
	if *edition != -1 {
		patchData["edition"] = *edition
	}
	if *description != "" {
		patchData["description"] = *description
	}
	if *genre != "" {
		patchData["genre"] = *genre
	}

	if len(patchData) == 0 {
		fmt.Println("No fields to update provided")
		os.Exit(1)
	}

	body, err := client.Patch(fmt.Sprintf("/books/%d", id), patchData)
	if err != nil {
		log.Fatalf("Error patching book: %v", err)
	}

	var patchedBook models.Book
	if err := json.Unmarshal(body, &patchedBook); err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	fmt.Printf("Successfully patched book #%d\n", patchedBook.ID)
}

func deleteBook(client *api.APIClient, args []string) {
	if len(args) < 1 {
		fmt.Println("Book ID is required")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid book ID")
		os.Exit(1)
	}

	err = client.Delete(fmt.Sprintf("/books/%d", id))
	if err != nil {
		log.Fatalf("Error deleting book: %v", err)
	}

	fmt.Printf("Deleted book #%d\n", id)
}
