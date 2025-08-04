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
)

func HandleCollectionCommand(client *api.APIClient, args []string) {
	if len(args) < 1 {
		printCollectionHelp()
		os.Exit(1)
	}

	switch args[0] {
	case "create":
		createCollection(client, args[1:])
	case "list":
		listCollections(client, args[1:])
	case "get":
		getCollection(client, args[1:])
	case "update":
		updateCollection(client, args[1:])
	case "patch":
		patchCollection(client, args[1:])
	case "delete":
		deleteCollection(client, args[1:])
	case "add-book":
		addBookToCollection(client, args[1:])
	case "remove-book":
		removeBookFromCollection(client, args[1:])
	case "list-books":
		listBooksInCollection(client, args[1:])
	case "help":
		printCollectionHelp()
	default:
		fmt.Printf("Unknown collection command: %s\n", args[0])
		printCollectionHelp()
		os.Exit(1)
	}
}

func printCollectionHelp() {

	fmt.Printf("%s", `Usage: bookmanager collection <command> [options]

Commands:
	create        Create a new collection
	list          List all collections with optional filters
	get           Get details of a specific collection
	update        Update a collection's information (requires all fields)
	patch					Partially update collection fields (only updates provided fields)
	delete        Remove a collection
	add-book      Add a book to a collection
	remove-book   Remove a book from a collection
	list-books    List books in a collection
	help          Show this help message

List Options:
	--where       SQL-like WHERE clause (e.g., "name LIKE '%Fantasy%'")
	--group-by    SQL-like GROUP BY clause (e.g., "SUBSTRING(name, 1, 1)")
	--order-by    SQL-like ORDER BY clause (e.g., "name DESC")
	--limit       Limit number of results
	--offset      Offset for pagination

Examples:
	bookmanager collection create --name "Fantasy Classics" --description "Classic fantasy books"
	bookmanager collection list --where "name LIKE '%Classics%'"
	bookmanager collection list --group-by "SUBSTRING(name, 1, 1)"`)
}

func createCollection(client *api.APIClient, args []string) {
	fs := flag.NewFlagSet("collection create", flag.ExitOnError)
	name := fs.String("name", "", "Collection name (required)")
	description := fs.String("description", "", "Collection description")
	fs.Parse(args)

	if *name == "" {
		fmt.Println("Name is required")
		fs.PrintDefaults()
		os.Exit(1)
	}

	collection := map[string]interface{}{
		"name":        *name,
		"description": *description,
	}

	body, err := client.Post("/collections", collection)
	if err != nil {
		log.Fatalf("Error creating collection: %v", err)
	}

	var createdCollection models.Collection
	if err := json.Unmarshal(body, &createdCollection); err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	fmt.Printf("Created collection #%d:%s\n", createdCollection.ID, createdCollection.Name)
}

func listCollections(client *api.APIClient, args []string) {
	fs := flag.NewFlagSet("collection list", flag.ExitOnError)
	where := fs.String("where", "", "SQL-like WHERE clause")
	groupBy := fs.String("group-by", "", "SQL-like GROUP BY clause")
	orderBy := fs.String("order-by", "", "SQL-like ORDER BY clause")
	limit := fs.Int("limit", 0, "Limit number of results")
	offset := fs.Int("offset", 0, "Offset for pagination")
	fs.Parse(args)

	params := client.BuildQueryParams(*where, *groupBy, *orderBy, *limit, *offset)

	body, err := client.Get("/v1/collections", params)
	if err != nil {
		log.Fatalf("Error listing collections: %v", err)
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

		fmt.Printf("Collections grouped by %s:\n", *groupBy)
		for group, count := range result.Groups {
			fmt.Printf("- %s: %d\n", group, count)
		}
	} else {
		var result struct {
			Collections []models.Collection `json:"collections"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatalf("Error parsing response: %v", err)
		}

		if len(result.Collections) == 0 {
			fmt.Println("No collections found")
			return
		}

		for _, collection := range result.Collections {
			fmt.Printf("%d: %s\n", collection.ID, collection.Name)
			if collection.Description != "" {
				fmt.Printf("   Description: %s\n", collection.Description)
			}
			fmt.Println()
		}
	}
}

func getCollection(client *api.APIClient, args []string) {
	if len(args) < 1 {
		fmt.Println("Collection ID is required")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid collection ID")
		os.Exit(1)
	}

	body, err := client.Get(fmt.Sprintf("/v1/collections/%d", id), nil)
	if err != nil {
		log.Fatalf("Error getting collection: %v", err)
	}

	var collection models.Collection
	if err := json.Unmarshal(body, &collection); err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	fmt.Printf("Collection #%d\n", collection.ID)
	fmt.Printf("Name: %s\n", collection.Name)
	if collection.Description != "" {
		fmt.Printf("Description: %s\n", collection.Description)
	}
}

func updateCollection(client *api.APIClient, args []string) {
	fs := flag.NewFlagSet("collection update", flag.ExitOnError)
	name := fs.String("name", "", "Collection name")
	description := fs.String("description", "", "Collection description")

	if len(args) < 1 {
		fmt.Println("Collection ID is required")
		fs.PrintDefaults()
		os.Exit(1)
	}

	fs.Parse(args[1:])

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid collection ID")
		os.Exit(1)
	}

	currentBody, err := client.Get(fmt.Sprintf("/collections/%d", id), nil)
	if err != nil {
		log.Fatalf("Error getting current collection: %v", err)
	}

	var currentCollection models.Collection
	if err := json.Unmarshal(currentBody, &currentCollection); err != nil {
		log.Fatalf("Error parsing current collection: %v", err)
	}

	updateData := make(map[string]interface{})
	if *name != "" {
		updateData["name"] = *name
	} else {
		updateData["name"] = currentCollection.Name
	}
	if *description != "" {
		updateData["description"] = *description
	} else {
		updateData["description"] = currentCollection.Description
	}

	body, err := client.Put(fmt.Sprintf("/collections/%d", id), updateData)
	if err != nil {
		log.Fatalf("Error updating collection: %v", err)
	}

	var updatedCollection models.Collection
	if err := json.Unmarshal(body, &updatedCollection); err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	fmt.Printf("Updated collection #%d: %s\n", updatedCollection.ID, updatedCollection.Name)
}

func patchCollection(client *api.APIClient, args []string) {
	fs := flag.NewFlagSet("collection patch", flag.ExitOnError)
	name := fs.String("name", "", "Update collection name (optional)")
	description := fs.String("description", "", "Update collection description (optional)")

	if len(args) < 1 {
		fmt.Println("Collection ID is required")
		fs.PrintDefaults()
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid collection ID")
		os.Exit(1)
	}

	fs.Parse(args[1:])

	patchData := make(map[string]interface{})

	if *name != "" {
		patchData["name"] = *name
	}
	if *description != "" {
		patchData["description"] = *description
	}

	if len(patchData) == 0 {
		fmt.Println("No fields to update provided")
		os.Exit(1)
	}

	body, err := client.Patch(fmt.Sprintf("/collections/%d", id), patchData)
	if err != nil {
		log.Fatalf("Error patching collection: %v", err)
	}

	var patchedCollection models.Collection
	if err := json.Unmarshal(body, &patchedCollection); err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	fmt.Printf("Successfully patched collection #%d\n", patchedCollection.ID)
}

func deleteCollection(client *api.APIClient, args []string) {
	if len(args) < 1 {
		fmt.Println("Collection ID is required")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid collection ID")
		os.Exit(1)
	}

	err = client.Delete(fmt.Sprintf("/collections/%d", id))
	if err != nil {
		log.Fatalf("Error deleting collection: %v", err)
	}

	fmt.Printf("Deleted collection #%d\n", id)
}

func addBookToCollection(client *api.APIClient, args []string) {
	if len(args) < 2 {
		fmt.Println("Collection ID and Book ID are required")
		os.Exit(1)
	}

	collectionID, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid collection ID")
		os.Exit(1)
	}

	bookID, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Invalid book ID")
		os.Exit(1)
	}

	req := map[string]interface{}{
		"book_id": bookID,
	}
	_, err = client.Post(fmt.Sprintf("/collections-books/%d", collectionID), req)
	if err != nil {
		log.Fatalf("Error adding book to collection: %v", err)
	}

	fmt.Printf("Added book #%d to collection #%d\n", bookID, collectionID)
}

func removeBookFromCollection(client *api.APIClient, args []string) {
	if len(args) < 2 {
		fmt.Println("Collection ID and Book ID are required")
		os.Exit(1)
	}

	collectionID, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid collection ID")
		os.Exit(1)
	}

	bookID, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Invalid book ID")
		os.Exit(1)
	}

	err = client.Delete(fmt.Sprintf("/collections-books/%d/%d", collectionID, bookID))
	if err != nil {
		log.Fatalf("Error removing book from collection: %v", err)
	}

	fmt.Printf("Removed book #%d from collection #%d\n", bookID, collectionID)
}

func listBooksInCollection(client *api.APIClient, args []string) {
	if len(args) < 1 {
		fmt.Println("Collection ID is required")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid collection ID")
		os.Exit(1)
	}

	body, err := client.Get(fmt.Sprintf("/v1/collections-books/%d", id), nil)
	if err != nil {
		log.Fatalf("Error listing books in collection: %v", err)

	}

	var result struct {
		Books []models.Book `json:"books"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	if len(result.Books) == 0 {
		fmt.Printf("No books found in collection #%d\n", id)
		return
	}

	fmt.Printf("Books in collection #%d:\n", id)
	for _, book := range result.Books {
		fmt.Printf("- %s by %s (%s)\n", book.Title, book.Author, book.PublishedDate)
	}
}
