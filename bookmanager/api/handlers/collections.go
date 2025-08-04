package handlers

import (
	"bookmanager/api/db"
	"bookmanager/api/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type CollectionHandler struct {
	db *db.CollectionDB
}

func NewCollectionHandler(db *db.CollectionDB) *CollectionHandler {
	return &CollectionHandler{db: db}
}

func (h *CollectionHandler) HandleCollections(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listCollections(w, r)
	case http.MethodPost:
		h.createCollection(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CollectionHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/api/v1/collections/"):])
	if err != nil {
		http.Error(w, "Invalid collection ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getCollection(w, r, id)
	case http.MethodPut:
		h.updateCollection(w, r, id)
	case http.MethodPatch:
		h.patchCollection(w, r, id)
	case http.MethodDelete:
		h.deleteCollection(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func (h *CollectionHandler) HandleCollectionBooks(w http.ResponseWriter, r *http.Request) {

	basePath := "/api/v1/collections-books/"
	if !strings.HasPrefix(r.URL.Path, basePath) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, basePath)
	parts := strings.Split(path, "/")
	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, "Collection ID required", http.StatusBadRequest)
		return
	}

	collectionID, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "Invalid collection ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.listBooksInCollection(w, r, collectionID)
	case http.MethodPost:
		h.addBookToCollection(w, r, collectionID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CollectionHandler) HandleCollectionBook(w http.ResponseWriter, r *http.Request) {

	basePath := "/api/v1/collections-books/"
	if !strings.HasPrefix(r.URL.Path, basePath) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, basePath)
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid URL path format: expected /collections-books/{collection_id}/{book_id}", http.StatusBadRequest)
		return
	}

	collectionID, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "Invalid collection ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	bookID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "Invalid book ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete {
		h.removeBookFromCollection(w, r, collectionID, bookID)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CollectionHandler) createCollection(w http.ResponseWriter, r *http.Request) {
	var collectionReq models.CollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&collectionReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := collectionReq.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection, err := h.db.CreateCollection(&collectionReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(collection)
}

func (h *CollectionHandler) listCollections(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	where := query.Get("where")
	groupBy := query.Get("group_by")
	orderBy := query.Get("order_by")
	limit := query.Get("limit")
	offset := query.Get("offset")

	collections, err := h.db.ListCollections(where, groupBy, orderBy, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if groupBy != "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"groups": collections})
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"collections": collections})
	}
}

func (h *CollectionHandler) getCollection(w http.ResponseWriter, r *http.Request, id int) {
	collection, err := h.db.GetCollection(id)
	if err != nil {
		if err.Error() == "collection not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collection)
}

func (h *CollectionHandler) updateCollection(w http.ResponseWriter, r *http.Request, id int) {
	var collectionReq models.CollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&collectionReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if collectionReq.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	collection, err := h.db.UpdateCollection(id, &collectionReq)
	if err != nil {
		if err.Error() == "collection not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collection)
}

func (h *CollectionHandler) patchCollection(w http.ResponseWriter, r *http.Request, id int) {
	var patch models.CollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	collection, err := h.db.PatchCollection(id, &patch)
	if err != nil {
		if err.Error() == "collection not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collection)
}

func (h *CollectionHandler) deleteCollection(w http.ResponseWriter, r *http.Request, id int) {
	err := h.db.DeleteCollection(id)
	if err != nil {
		if err.Error() == "collection not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CollectionHandler) addBookToCollection(w http.ResponseWriter, r *http.Request, collectionID int) {
	var req struct {
		BookID int `json:"book_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.BookID == 0 {
		http.Error(w, "book_id is required", http.StatusBadRequest)
		return
	}

	err := h.db.AddBookToCollection(collectionID, req.BookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&req)
	w.WriteHeader(http.StatusCreated)
}

func (h *CollectionHandler) removeBookFromCollection(w http.ResponseWriter, r *http.Request, collectionID, bookID int) {
	err := h.db.RemoveBookFromCollection(collectionID, bookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CollectionHandler) listBooksInCollection(w http.ResponseWriter, r *http.Request, collectionID int) {
	books, err := h.db.ListBooksInCollection(collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"books": books})
}

func (h *CollectionHandler) HandleCollectionBooksRoutes(w http.ResponseWriter, r *http.Request) {
	cleanPath := strings.TrimPrefix(r.URL.Path, "/api/v1/collections-books")
	cleanPath = strings.TrimPrefix(cleanPath, "/")
	parts := strings.Split(cleanPath, "/")

	if len(parts) >= 2 && parts[1] != "" {
		h.HandleCollectionBook(w, r)
	} else if len(parts) == 1 && parts[0] != "" {
		h.HandleCollectionBooks(w, r)
	} else {
		http.Error(w, "Invalid path", http.StatusBadRequest)
	}
}
