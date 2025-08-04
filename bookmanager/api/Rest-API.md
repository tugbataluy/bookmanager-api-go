
# BookManager API Documentation

## Table of Contents

- [Book Requests](#book-requests)
    - [Create Book Record](#create-book-record)
    - [Get All Book Records](#get-all-book-records)
    - [Get Specific Book Record](#get-specific-book-record)
    - [Delete Book](#delete-book)
    - [Update Book (Full)](#update-book-full)
    - [Update Book (Partial)](#update-book-partial)
- [Collection Requests](#collection-requests)
    - [Create Collection](#create-collection)
    - [Get All Collections](#get-all-collections)
    - [Get Specific Collection](#get-specific-collection)
    - [Delete Collection](#delete-collection)
    - [Update Collection (Full)](#update-collection-full)
    - [Update Collection (Partial)](#update-collection-partial)
- [Collection Books Requests](#collection-books-requests)
    - [Add Book to Collection](#add-book-to-collection)
    - [List Books in a Collection](#list-books-in-a-collection)
    - [Delete Book from a Collection](#delete-book-from-a-collection)

---
## Status Codes

The following HTTP status codes are used throughout the BookManager API:

| Status Code | Meaning                | Where Used                                                                                      |
|-------------|------------------------|-------------------------------------------------------------------------------------------------|
| 200 OK      | Success                | All successful GET, PUT, PATCH requests (books, collections, collection-books)                  |
| 201 Created | Resource created       | Creating books, collections, adding a book to a collection                                      |
| 204 No Content | Resource deleted    | Deleting books, collections, removing a book from a collection                                  |
| 400 Bad Request | Invalid input      | Invalid input or request for create, update, patch, or handler endpoints                        |
| 404 Not Found   | Resource not found | When a requested book, collection, or collection-book does not exist                            |
| 405 Method Not Allowed | Not allowed | When an unsupported HTTP method is used on an endpoint                                          |
| 500 Internal Server Error | Server error | Any unexpected server error during create, list, get, update, patch, or delete operations   |


## Book Requests

### Create Book Record

- **Endpoint:** `POST /api/v1/books`
- **Example URL:** `http://localhost:8080/api/v1/books`
- **Request Body:**
    ```json
    {
        "title": "Dune Messiah",
        "author": "Frank Herbert",
        "published_date": "1969-01-01",
        "edition": 1,
        "description": "The second book in the Dune series",
        "genre": "Science Fiction"
    }
    ```
- **Example cURL:**
    ```sh
    curl -X POST http://localhost:8080/api/v1/books \
        -H "Content-Type: application/json" \
        -d '{"Title": "Dune", "Author": "Frank Herbert", "published_date": "1965-08-01", "Edition": 1, "Description": "The first book in the Dune series", "Genre": "Science Fiction"}'
    ```
- **Response:**
    ```json
    {
        "id": 1,
        "title": "Dune",
        "author": "Frank Herbert",
        "published_date": "1965-08-01T00:00:00Z",
        "edition": 1,
        "description": "The first book in the Dune series",
        "genre": "Science Fiction",
        "created_at": "...",
        "updated_at": "..."
    }
    ```

---

### Get All Book Records

- **Endpoint:** `GET /api/v1/books`
- **Example URL:** `http://localhost:8080/api/v1/books`
- **Request Body:** None
- **Example cURL:**
    ```sh
    curl -X GET http://localhost:8080/api/v1/books
    ```
- **Response:**
    ```json
    {
        "books": [
            {
                "id": 1,
                "title": "Dune",
                "author": "Frank Herbert",
                "published_date": "1965-08-01T00:00:00Z",
                "edition": 1,
                "description": "The first book in the Dune series",
                "genre": "Science Fiction",
                "created_at": "...",
                "updated_at": "..."
            },
            {
                "id": 2,
                "title": "Dune Messiah",
                "author": "Frank Herbert",
                "published_date": "1969-01-01T00:00:00Z",
                "edition": 1,
                "description": "The second book in the Dune series",
                "genre": "Science Fiction",
                "created_at": "...",
                "updated_at": "..."
            }
        ]
    }
    ```

---

### Get Specific Book Record

- **Endpoint:** `GET /api/v1/books/{book_id}`
- **Example URL:** `http://localhost:8080/api/v1/books/3`
- **Request Body:** None
- **Example cURL:**
    ```sh
    curl -X GET http://localhost:8080/api/v1/books/3
    ```
- **Response:**
    ```json
    {
        "id": 3,
        "title": "Children of Dune",
        "author": "Frank Herbert",
        "published_date": "1976-01-01T00:00:00Z",
        "edition": 1,
        "description": "The third book in the Dune series",
        "genre": "Science Fiction",
        "created_at": "...",
        "updated_at": "..."
    }
    ```

---

### Delete Book

- **Endpoint:** `DELETE /api/v1/books/{book_id}`
- **Example URL:** `http://localhost:8080/api/v1/books/7`
- **Request Body:** None
- **Example cURL:**
    ```sh
    curl -X DELETE http://localhost:8080/api/v1/books/7
    ```
- **Response:** None (if successful)

---

### Update Book (Full)

- **Endpoint:** `PUT /api/v1/books/{book_id}`
- **Example URL:** `http://localhost:8080/api/v1/books/5`
- **Request Body:**
    ```json
    {
        "title": "Heretics of Dune (Updated)",
        "author": "Frank Herbert",
        "published_date": "1984-01-01",
        "edition": 2,
        "description": "The fifth book in the Dune series",
        "genre": "Science Fiction"
    }
    ```
- **Example cURL:**
    ```sh
    curl -X PUT http://localhost:8080/api/v1/books/5 \
        -H "Content-Type: application/json" \
        -d '{"Title": "Heretics of Dune (Updated)", "Author": "Frank Herbert", "published_date": "1984-01-01", "Edition": 2, "Description": "The fifth book in the Dune series", "Genre": "Science Fiction"}'
    ```
- **Response:**
    ```json
    {
        "id": 5,
        "title": "Heretics of Dune (Updated)",
        "author": "Frank Herbert",
        "published_date": "1984-01-01T00:00:00Z",
        "edition": 2,
        "description": "The fifth book in the Dune series",
        "genre": "Science Fiction",
        "created_at": "...",
        "updated_at": "..."
    }
    ```

---

### Update Book (Partial)

- **Endpoint:** `PATCH /api/v1/books/{book_id}`
- **Example URL:** `http://localhost:8080/api/v1/books/5`
- **Request Body:** (any subset of fields)
    ```json
    {
        "description": "The fifth book in the Dune series (My fav)",
        "genre": "Science Fiction"
    }
    ```
- **Example cURL:**
    ```sh
    curl -X PATCH http://localhost:8080/api/v1/books/5 \
        -H "Content-Type: application/json" \
        -d '{"Description": "The fifth book in the Dune series (My fav)", "Genre": "Science Fiction"}'
    ```
- **Response:**
    ```json
    {
        "id": 5,
        "title": "Heretics of Dune (Updated)",
        "author": "Frank Herbert",
        "published_date": "1984-01-01T00:00:00Z",
        "edition": 2,
        "description": "The fifth book in the Dune series (My fav)",
        "genre": "Science Fiction",
        "created_at": "...",
        "updated_at": "..."
    }
    ```

---

## Collection Requests

### Create Collection

- **Endpoint:** `POST /api/v1/collections`
- **Example URL:** `http://localhost:8080/api/v1/collections`
- **Request Body:**
    ```json
    {
        "name": "Fantasy Novels",
        "description": "A collection of fantasy books."
    }
    ```
- **Example cURL:**
    ```sh
    curl -X POST http://localhost:8080/api/v1/collections \
        -H "Content-Type: application/json" \
        -d '{"name": "Fantasy Novels", "description": "A collection of fantasy books."}'
    ```
- **Response:**
    ```json
    {
        "id": 1,
        "name": "Fantasy Novels",
        "description": "A collection of fantasy books.",
        "created_at": "...",
        "updated_at": "..."
    }
    ```

---

### Get All Collections

- **Endpoint:** `GET /api/v1/collections`
- **Example URL:** `http://localhost:8080/api/v1/collections`
- **Request Body:** None
- **Example cURL:**
    ```sh
    curl -X GET http://localhost:8080/api/v1/collections
    ```
- **Response:**
    ```json
    {
        "collections": [
            {
                "id": 2,
                "name": "Dummy",
                "description": "",
                "created_at": "...",
                "updated_at": "..."
            },
            {
                "id": 1,
                "name": "Fantasy Novels",
                "description": "A collection of fantasy books.",
                "created_at": "...",
                "updated_at": "..."
            }
        ]
    }
    ```

---

### Get Specific Collection

- **Endpoint:** `GET /api/v1/collections/{collection_id}`
- **Example URL:** `http://localhost:8080/api/v1/collections/1`
- **Request Body:** None
- **Example cURL:**
    ```sh
    curl -X GET http://localhost:8080/api/v1/collections/1
    ```
- **Response:**
    ```json
    {
        "id": 1,
        "name": "Fantasy Novels",
        "description": "A collection of fantasy books.",
        "created_at": "...",
        "updated_at": "..."
    }
    ```

---

### Delete Collection

- **Endpoint:** `DELETE /api/v1/collections/{collection_id}`
- **Example URL:** `http://localhost:8080/api/v1/collections/2`
- **Request Body:** None
- **Example cURL:**
    ```sh
    curl -X DELETE http://localhost:8080/api/v1/collections/2
    ```
- **Response:** None (if successful)

---

### Update Collection (Full)

- **Endpoint:** `PUT /api/v1/collections/{collection_id}`
- **Example URL:** `http://localhost:8080/api/v1/collections/1`
- **Request Body:**
    ```json
    {
        "name": "Science Fiction Novels",
        "description": "A collection of science fiction books."
    }
    ```
- **Example cURL:**
    ```sh
    curl -X PUT http://localhost:8080/api/v1/collections/1 \
        -H "Content-Type: application/json" \
        -d '{"name": "Science Fiction Novels", "description": "A collection of science fiction books."}'
    ```
- **Response:**
    ```json
    {
        "id": 1,
        "name": "Science Fiction Novels",
        "description": "A collection of science fiction books.",
        "created_at": "...",
        "updated_at": "..."
    }
    ```

---

### Update Collection (Partial)

- **Endpoint:** `PATCH /api/v1/collections/{collection_id}`
- **Example URL:** `http://localhost:8080/api/v1/collections/1`
- **Request Body:** (any subset of fields)
    ```json
    {
        "description": "A collection of sci-fi books."
    }
    ```
- **Example cURL:**
    ```sh
    curl -X PATCH http://localhost:8080/api/v1/collections/1 \
        -H "Content-Type: application/json" \
        -d '{ "description": "A collection of sci-fi books."}'
    ```
- **Response:**
    ```json
    {
        "id": 1,
        "name": "Science Fiction Novels",
        "description": "A collection of sci-fi books.",
        "created_at": "...",
        "updated_at": "..."
    }
    ```

---

## Collection Books Requests

### Add Book to Collection

- **Endpoint:** `POST /api/v1/collections-books/{collection_id}`
- **Example URL:** `http://localhost:8080/api/v1/collections-books/1`
- **Request Body:**
    ```json
    {
        "book_id": 1
    }
    ```
- **Example cURL:**
    ```sh
    curl -X POST http://localhost:8080/api/v1/collections-books/1 \
        -H "Content-Type: application/json" \
        -d '{ "book_id": 1 }'
    ```
- **Response:**
    ```json
    {
        "book_id": 1
    }
    ```

---

### List Books in a Collection

- **Endpoint:** `GET /api/v1/collections-books/{collection_id}`
- **Example URL:** `http://localhost:8080/api/v1/collections-books/1`
- **Request Body:** None
- **Example cURL:**
    ```sh
    curl -X GET http://localhost:8080/api/v1/collections-books/1
    ```
- **Response:**
    ```json
    {
        "books": [
            {
                "id": 1,
                "title": "Dune",
                "author": "Frank Herbert",
                "published_date": "1965-08-01T00:00:00Z",
                "edition": 1,
                "description": "The first book in the Dune series",
                "genre": "Science Fiction",
                "created_at": "...",
                "updated_at": "..."
            },
            {
                "id": 2,
                "title": "Dune Messiah",
                "author": "Frank Herbert",
                "published_date": "1969-01-01T00:00:00Z",
                "edition": 1,
                "description": "The second book in the Dune series",
                "genre": "Science Fiction",
                "created_at": "...",
                "updated_at": "..."
            }
        ]
    }
    ```

---

### Delete Book from a Collection

- **Endpoint:** `DELETE /api/v1/collections-books/{collection_id}/{book_id}`
- **Example URL:** `http://localhost:8080/api/v1/collections-books/1/2`
- **Request Body:** None
- **Example cURL:**
    ```sh
    curl -X DELETE http://localhost:8080/api/v1/collections-books/1/2
    ```
- **Response:** None (if successful)

---

