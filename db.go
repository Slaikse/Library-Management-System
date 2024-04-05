package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/boltdb/bolt"
)

var db *bolt.DB
var err error

// Book структура, описывающая книгу
type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
}

// Инициализация базы данных
func initDB() {
	db, err = bolt.Open("library.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Books"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

// Добавление новой книги
func addBook(book Book) error {
	err := addGenre(book.Genre) // Добавляем жанр перед добавлением книги
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Books"))
		id, _ := b.NextSequence()
		book.ID = fmt.Sprintf("%d", id)

		buf, err := json.Marshal(book)
		if err != nil {
			return err
		}
		return b.Put([]byte(book.ID), buf)
	})
}

// Получение списка всех книг
func getAllBooks() ([]Book, error) {
	var books []Book
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Books"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var book Book
			err := json.Unmarshal(v, &book)
			if err != nil {
				return err
			}
			books = append(books, book)
		}
		return nil
	})
	return books, err
}

// searchBooks осуществляет поиск книг по ключевому слову в названии и авторе
func searchBooks(query string) ([]Book, error) {
	var matchedBooks []Book
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Books"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var book Book
			if err := json.Unmarshal(v, &book); err != nil {
				return err
			}
			if strings.Contains(strings.ToLower(book.Title), strings.ToLower(query)) ||
				strings.Contains(strings.ToLower(book.Author), strings.ToLower(query)) {
				matchedBooks = append(matchedBooks, book)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matchedBooks, nil
}

// Удаление книги по ID
func deleteBook(bookID string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Books"))
		return b.Delete([]byte(bookID))
	})
}

// Вспомогательная функция для проверки содержания подстроки
func contains(source, toFind string) bool {
	return strings.Contains(strings.ToLower(source), strings.ToLower(toFind))
}

// Добавление нового жанра, если он уникален
func addGenre(genreName string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Genres"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		// Проверяем, существует ли жанр
		exists := b.Get([]byte(genreName))
		if exists == nil {
			return b.Put([]byte(genreName), []byte(genreName))
		}
		return nil
	})
}

// Получение списка всех жанров
func getGenres() ([]string, error) {
	var genres []string
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Genres"))
		if b == nil {
			return fmt.Errorf("Genres bucket does not exist")
		}
		return b.ForEach(func(k, v []byte) error {
			genres = append(genres, string(v))
			return nil
		})
	})
	return genres, err
}

// Получение списка книг по жанру
func getBooksByGenre(genre string) ([]Book, error) {
	var booksByGenre []Book
	books, err := getAllBooks()
	if err != nil {
		return nil, err
	}

	for _, book := range books {
		if strings.EqualFold(book.Genre, genre) {
			booksByGenre = append(booksByGenre, book)
		}
	}
	return booksByGenre, nil
}

// Получение книги по ID
func getBookByID(bookID string) (*Book, error) {
	var book *Book
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Books"))
		bookData := b.Get([]byte(bookID))
		if bookData == nil {
			return fmt.Errorf("book not found")
		}
		if err := json.Unmarshal(bookData, &book); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return book, nil
}
