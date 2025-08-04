# BookManager CLI

A modern command-line interface for managing your books and collections via the BookManager API.

---

## Building

```sh
# Build the API server
go build -o bookmanagar_api ./api

# Build the CLI
go build -o bookmanager bin/cli ./cmd/bookmanager
```

---

## Usage

```sh
./bookmanager
```

```
Usage: bookmanager [global options] <command> [command options]

Global options:
    --api-url    URL of the API server (default: http://localhost:8080/api/v1)
    --verbose    Enable verbose output
    --version    Show version and exit
    --help       Show help

Commands:
    book        Manage books
    collection  Manage collections
    help        Shows this help message

Use 'bookmanager <command> --help' for more information about a command.
```

---

## Book Commands

### Command Reference

```sh
./bookmanager book
```

```
Usage: bookmanager book <command> [options]

Commands:
    create      Add a new book
    list        List books with optional filters
    get         Get details of a specific book
    update      Update a book's information (all fields required)
    patch       Partially update a book's information
    delete      Remove a book from the system
    help        Show this help message
```

#### List Options

- `--author`           Filter by author
- `--genre`            Filter by genre
- `--published-after`  Filter by publication date (after)
- `--published-before` Filter by publication date (before)
- `--where`            SQL-like WHERE clause (e.g., `"title LIKE '%Hobbit%' AND edition > 1"`)
- `--order-by`         SQL-like ORDER BY clause (e.g., `"published_date DESC"`)
- `--group-by`         SQL-like GROUP BY clause (e.g., `"genre"`)
- `--limit`            Limit number of results
- `--offset`           Offset for pagination

---

### Examples

#### Create a Book

```sh
bookmanager book create --title "The Go Programming Language" \
    --author "Alan A. A. Donovan & Brian W. Kernighan" \
    --published-date "2015-10-26" \
    --edition 1 \
    --description "Authoritative resource for learning Go" \
    --genre "Programming"
```
**Output:**
```
Created book 8: The Go Programming Language by Alan A. A. Donovan & Brian W. Kernighan
```

#### List All Books

```sh
bookmanager book list
```
**Output:**
```
6: Chapterhouse: Dune by Frank Herbert (1985-01-01T00:00:00Z)
   Genre: Science Fiction
   Description: The sixth book in the Dune series
...
8: The Go Programming Language by Alan A. A. Donovan & Brian W. Kernighan (2015-10-26T00:00:00Z)
   Genre: Programming
   Description: Authoritative resource for learning Go
```

#### Filter Books by Genre

```sh
bookmanager book list --genre "Software Engineering"
```
**Output:**
```
9: Clean Code: A Handbook of Agile Software Craftsmanship by Robert C. Martin (2008-08-01T00:00:00Z)
   Genre: Software Engineering
   Description: Principles for writing maintainable code
...
```

#### Filter Books by Author

```sh
bookmanager book list --author "Frank Herbert"
```
**Output:**
```
6: Chapterhouse: Dune by Frank Herbert (1985-01-01T00:00:00Z)
   Genre: Science Fiction
   Description: The sixth book in the Dune series
...
```

#### Group Books by Genre

```sh
bookmanager book list --group-by genre
```
**Output:**
```
Books grouped by genre:
- Computer Science: 1
- Programming: 1
- Science Fiction: 6
- Software Engineering: 2
```

#### Combine List Commands

```sh
bookmanager book list --limit 4 --order-by "title DESC"
```
**Output:**
```
8: The Go Programming Language by Alan A. A. Donovan & Brian W. Kernighan (2015-10-26T00:00:00Z)
   Genre: Programming
   Description: Authoritative resource for learning Go
...
```

#### Get Book Details

```sh
bookmanager book get 3
```
**Output:**
```
Book #3
Title: Children of Dune
Author: Frank Herbert
Published Date: 1976-01-01T00:00:00Z
Edition: 1
Description: The third book in the Dune series
Genre: Science Fiction
```

#### Delete a Book

```sh
bookmanager book delete 13
```
**Output:**
```
Deleted book #13
```

#### Update a Book

```sh
bookmanager book update 10 --title "Design Patterns: Elements of Reusable Object-Oriented Software" \
    --author "Erich Gamma, Richard Helm, Ralph Johnson, John Vlissides" \
    --published-date "1994-10-31" \
    --edition 2 \
    --description "Classic solutions to common design problems" \
    --genre "Software Engineering"
```
**Output:**
```
Updated book #10: Design Patterns: Elements of Reusable Object-Oriented Software by Erich Gamma, Richard Helm, Ralph Johnson, John Vlissides
```

#### Patch a Book

```sh
bookmanager book patch 10 --author "Erich Gamma, Richard Helm, Ralph Johnson, Me"
```
**Output:**
```
Successfully patched book #10
```

---

## Collection Commands

### Command Reference

```sh
./bookmanager collection
```

```
Usage: bookmanager collection <command> [options]

Commands:
    create        Create a new collection
    list          List all collections with optional filters
    get           Get details of a specific collection
    update        Update a collection's information (requires all fields)
    patch         Partially update collection fields
    delete        Remove a collection
    add-book      Add a book to a collection
    remove-book   Remove a book from a collection
    list-books    List books in a collection
    help          Show this help message
```

#### List Options

- `--where`       SQL-like WHERE clause (e.g., `"name LIKE '%Fantasy%'"`)
- `--group-by`    SQL-like GROUP BY clause (e.g., `"SUBSTRING(name, 1, 1)"`)
- `--order-by`    SQL-like ORDER BY clause (e.g., `"name DESC"`)
- `--limit`       Limit number of results
- `--offset`      Offset for pagination

---

### Examples

#### List All Collections

```sh
bookmanager collection list
```
**Output:**
```
1: Science Fiction Novels
   Description: A collection of sci-fi books.
```

#### Create a Collection

```sh
bookmanager collection create --name "Computer Science Classics"
```
**Output:**
```
Created collection #3:Computer Science Classics
```

#### Add Book to Collection

```sh
bookmanager collection add-book 3 10
```
**Output:**
```
Added book #10 to collection #3
```

#### List Books in a Collection

```sh
bookmanager collection list-books 3
```
**Output:**
```
Books in collection #3:
- Design Patterns: Elements of Reusable Object-Oriented Software by Erich Gamma, Richard Helm, Ralph Johnson (1994-10-31T00:00:00Z)
```

#### Remove Book from Collection

```sh
bookmanager collection remove-book 3 10
```
**Output:**
```
Removed book #10 from collection #3
```

#### Group Collections

```sh
bookmanager collection list --group-by "description"
```
**Output:**
```
Collections grouped by description:
- : 1
- A collection of sci-fi books.: 1
```

#### Update a Collection

```sh
bookmanager collection update 3 --name "Foundation Of CS" --description "Computer science books for enthusiast"
```
**Output:**
```
Updated collection #3: Foundation Of CS
```

#### Patch a Collection

```sh
bookmanager collection patch 3 --description "Books for CS enthusiasts"
```
**Output:**
```
Successfully patched collection #3
```

#### Delete a Collection

```sh
bookmanager --verbose collection delete 3
```
**Output:**
```
DELETE http://localhost:8080/api/v1/collections/3
Deleted collection #3
```

---

## Help

For more information on any command, use:

```sh
bookmanager <command> --help
```

---

Enjoy managing your books and collections from the command line!
