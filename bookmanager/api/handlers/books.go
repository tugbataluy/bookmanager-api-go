package handlers

import (
	"bookmanager/api/db"
	"bookmanager/api/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BookHandler struct {
	db *db.BookDB
}

func NewBookHandler(db *db.BookDB) *BookHandler {
	return &BookHandler{db: db}
}

func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listBooks(w, r)
	case http.MethodPost:
		h.createBook(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BookHandler) HandleBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/api/v1/books/"):])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getBook(w, r, id)
	case http.MethodPut:
		h.updateBook(w, r, id)
	case http.MethodPatch:
		h.patchBook(w, r, id)
	case http.MethodDelete:
		h.deleteBook(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	var bookReq models.BookRequest
	if err := json.NewDecoder(r.Body).Decode(&bookReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := bookReq.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := h.db.CreateBook(&bookReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) listBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	where := query.Get("where")
	groupBy := query.Get("group_by")
	orderBy := query.Get("order_by")
	limit := query.Get("limit")
	offset := query.Get("offset")

	// Extract additional book-specific filters
	author := query.Get("author")
	genre := query.Get("genre")
	publishedAfter := query.Get("published_after")
	publishedBefore := query.Get("published_before")

	// Build where clause combining manual where and filters
	var whereClauses []string
	if where != "" {
		whereClauses = append(whereClauses, "("+where+")")
	}
	if author != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("author ILIKE '%%%s%%'", author))
	}
	if genre != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("genre = '%s'", genre))
	}
	if publishedAfter != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("published_date >= '%s'", publishedAfter))
	}
	if publishedBefore != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("published_date <= '%s'", publishedBefore))
	}

	combinedWhere := strings.Join(whereClauses, " AND ")

	books, err := h.db.ListBooks(combinedWhere, groupBy, orderBy, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if groupBy != "" {
		// Handle grouped response
		if groups, ok := books.(map[string]int); ok {
			json.NewEncoder(w).Encode(map[string]interface{}{"groups": groups})
		} else {
			http.Error(w, "unexpected group response type", http.StatusInternalServerError)
		}
	} else {
		// Handle normal book list response
		if bookList, ok := books.([]models.Book); ok {
			json.NewEncoder(w).Encode(map[string]interface{}{"books": bookList})
		} else {
			http.Error(w, "unexpected book response type", http.StatusInternalServerError)
		}
	}
}

func (h *BookHandler) getBook(w http.ResponseWriter, r *http.Request, id int) {
	book, err := h.db.GetBook(id)
	if err != nil {
		if err.Error() == "book not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request, id int) {
	var bookReq models.BookRequest
	if err := json.NewDecoder(r.Body).Decode(&bookReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := bookReq.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := h.db.UpdateBook(id, &bookReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) patchBook(w http.ResponseWriter, r *http.Request, id int) {
	var patch models.BookRequest
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if patch.PublishedDate != "" {
		if _, err := time.Parse("2006-01-02", patch.PublishedDate); err != nil {
			http.Error(w, "invalid published_date format", http.StatusBadRequest)
			return
		}
	}

	book, err := h.db.PatchBook(id, &patch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request, id int) {
	err := h.db.DeleteBook(id)
	if err != nil {
		if err.Error() == "book not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
