package services

import (
	"fmt"

	m "github.com/vhladko/books/models"
	r "github.com/vhladko/books/repositories"
)

func GetBookByIsbn(isbn string) (m.Book, error) {
	book, err := r.GetBookByIsbn(isbn)
	if err == nil {
		return book, nil
	}

	fmt.Print("started from good reads")
	book, err = GetBookFromGoodreads(isbn)

	if err != nil {
		return m.Book{}, err
	}

	author, err := GetAuthorForGoodreads(book)
	book.AuthorId = author.Id
	r.AddBook(&book)

	return book, err
}

func GetAllBooks() []m.Book {
	books := r.GetAllBooks()

	return books;
}
