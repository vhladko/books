package repositories

import (
	m "github.com/vhladko/books/models"
)


func GetBookByIsbn(isbn string) (m.Book, error) {
	var book = m.Book{Isbn: isbn}

	err := DB.
		QueryRow(`select id, created_at, title, author_id, author_name, total_pages, description, preview_url, isbn from "public"."book" where "isbn"=$1`, isbn).
		Scan(&book.Id, &book.CreatedAt, &book.Title, &book.AuthorId, &book.AuthorName, &book.TotalPages, &book.Description, &book.PreviewUrl, &book.Isbn)

	if err != nil {
		return book, err
	} else {
		return book, nil
	}
}

func GetBookById(id string) (m.Book, error) {
	var book = m.Book{Id: id}

	err := DB.
		QueryRow(`select id, created_at, title, author_id, author_name, total_pages, description, preview_url, isbn from "public"."book" where "id"=$1`, id).
		Scan(&book.Id, &book.CreatedAt, &book.Title, &book.AuthorId, &book.AuthorName, &book.TotalPages, &book.Description, &book.PreviewUrl, &book.Isbn)

	if err != nil {
		return book, err
	} else {
		return book, nil
	}
}

func AddBook(book *m.Book) {
	insertBookQuery := `insert into "public"."book"("title", "total_pages", "description", "preview_url", "isbn", "author_id", "author_name") values($1, $2, $3, $4, $5, $6, $7) returning id, created_at`
	e := DB.QueryRow(insertBookQuery, book.Title, book.TotalPages, book.Description, book.PreviewUrl, book.Isbn, book.AuthorId, book.AuthorName).Scan(&book.Id, &book.CreatedAt)

	if e != nil {
		panic(e)
	}
}
