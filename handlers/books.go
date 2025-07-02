package handlers

import (
	"encoding/json"
	"library-api/models"
	"library-api/storage"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"
)

/*
Handlers are the glue between the web world (HTTP requests) and our application's logic (the storage layer). A handler's job is to:
Read and parse an incoming HTTP request.
Call the appropriate method(s) on the storage layer.
Write an HTTP response back to the client.
We'll use a lightweight and popular router called chi. Let's add it to our project:
*/

type BookHandler struct {
	Store *storage.Storage
}

func NewBookHandler(s *storage.Storage) *BookHandler {
	return &BookHandler{Store: s}
}

/*
	Crucial Explainer: Dependency Injection
Notice how our BookHandler doesn't create its own Storage.
Instead, it receives a *storage.Storage instance when it's created (NewBookHandler).
 This is a powerful pattern called Dependency Injection.
Why is this good? It decouples our components.
The handler doesn't care how the storage is created, only that it fulfills the contract of what a Storage can do.
This makes our code incredibly easy to test—we can "inject" a fake, or "mock," storage layer during tests instead of needing a real database.
*/

func (h *BookHandler) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Store.CreateBook(&book); err != nil {
		http.Error(w, "error while creating the book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) HandleGetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.Store.GetBooks()

	if err != nil {
		http.Error(w, "error while fetching books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "applikcation/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) HandleGetBookbyId(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid Book id", http.StatusBadRequest)
		return
	}

	book, err := h.Store.GetBookbyId(id)

	if err != nil {
		if err.Error() == "book not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "failed to retrieve the book", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) HandleUpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid value for book-id", http.StatusBadRequest)
		return
	}

	var book models.Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "invalid rfequest body", http.StatusBadRequest)
		return
	}

	if err := h.Store.UpdateBook(id, &book); err != nil {
		http.Error(w, "Failed to update book ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "book updated successfully"})
}

func (h *BookHandler) HandleDeleteBook(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid value for book-id", http.StatusBadRequest)
		return
	}

	if err := h.Store.DeleteBook(id); err != nil {
		http.Error(w, "error while delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
