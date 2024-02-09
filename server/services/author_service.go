package services

import (
	m "github.com/vhladko/books/models"
	r "github.com/vhladko/books/repositories"
)

func GetAuthorForGoodreads(book m.Book) (m.Author, error) {
	author, err := r.GetAuthorByName(book.AuthorName)

	if err == nil {
		return author, err
	}

	author, err = r.AddAuthor(book.AuthorName)

	return author, err

}
