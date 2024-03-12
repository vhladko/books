package repositories

import (
	m "github.com/vhladko/books/models"
)

func AddUser(user m.User) (m.User, error) {
	insertUserQuery := `insert into user_(email, password ) values($1, $2) returning id, created_at, email`
	err := DB.QueryRow(insertUserQuery, user.Email, user.Password).Scan(&user.Id, &user.CreatedAt, &user.Email)

	if err != nil {
		return user, err
	}

	return user, nil
}

func GetUserByEmail(email string) (m.User, error) {
	var user = m.User{}

	err := DB.
		QueryRow(`select id, created_at, email, username, password from user_ where email=$1`, email).
		Scan(&user.Id, &user.CreatedAt, &user.Email, &user.Username, &user.Password)

	if err != nil {
		return user, err
	} else {
		return user, nil
	}
}

func GetUserByUsername(username string) (m.User, error) {
	var user = m.User{}

	err := DB.
		QueryRow(`select id, created_at, email, username, password from user_ where username=$1`, username).
		Scan(&user.Id, &user.CreatedAt, &user.Email, &user.Username, &user.Password)

	if err != nil {
		return user, err
	} else {
		return user, nil
	}
}

func GetUserBooks(user m.User) []m.Book {
	var books []m.Book

	rows, err := DB.Query(`select b.id, b.isbn, b.created_at, b.title, b.author_id, b.total_pages, b.description, b.preview_url, b.author_name from user_book ub where user_id=$1 inner join book b on ub.book_id = b.id`, user.Id)


	if err != nil {
		return books
	}

	defer rows.Close()

	for rows.Next() {
		var book m.Book

		rows.Scan(&book.Id, &book.Isbn, &book.CreatedAt, &book.Title, &book.AuthorId, &book.TotalPages, &book.Description, &book.PreviewUrl, &book.AuthorName)

		books = append(books, book)
	}

	return books
}

func AddBookToUser(userId string, bookId string) error {
	insertUserQuery := `insert into user_book(user_id, book_id, status, total_pages_read ) values($1, $2, "none", 0)`
	DB.QueryRow(insertUserQuery, userId, bookId)

	return nil
}
