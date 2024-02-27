package repositories

import (
	m "github.com/vhladko/books/models"
)

func GetAuthorByName(name string) (m.Author, error) {
	var author = m.Author{}
	err := DB.
		QueryRow(`select id from author where name=$1`, name).
		Scan(&author.Id)

	if err != nil {
		return author, err
	} else {
		return author, nil
	}
}

func AddAuthor(name string) (m.Author, error) {
	var author = m.Author{}
	insertBookQuery := `insert into author(name) values($1) returning id, created_at, name`
	err := DB.QueryRow(insertBookQuery, name).Scan(&author.Id, &author.CreatedAt, &author.Name)

	if err != nil {
		return author, err
	}

	return author, nil

}
