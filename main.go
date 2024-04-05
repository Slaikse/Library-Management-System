package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	// Маршруты
	r.HandleFunc("/books", getBooksHandler).Methods("GET")
	r.HandleFunc("/book", addBookHandler).Methods("POST")
	r.HandleFunc("/book/{id}", getBookHandler).Methods("GET")
	r.HandleFunc("/book/{id}", deleteBookHandler).Methods("DELETE")
	r.HandleFunc("/genres", getGenresHandler).Methods("GET")
	r.HandleFunc("/books/genre/{genre}", getBooksByGenreHandler).Methods("GET")

	// Запуск сервера
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getBooksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := getAllBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}

func addBookHandler(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := addBook(book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(book) // Возвращаем добавленную книгу
}

func getBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]
	book, err := getBookByID(bookID)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]
	if err := deleteBook(bookID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK) // Успешное удаление
}

func getGenresHandler(w http.ResponseWriter, r *http.Request) {
	genres, err := getGenres()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(genres)
}

func getBooksByGenreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	genre := vars["genre"]
	books, err := getBooksByGenre(genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}
