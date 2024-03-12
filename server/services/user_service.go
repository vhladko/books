package services

import (
	"github.com/vhladko/books/jwt"
	m "github.com/vhladko/books/models"
	r "github.com/vhladko/books/repositories"
	"golang.org/x/crypto/bcrypt"
)

func AddUser(email string, password string) (m.User, error) {
	user := m.User{Email: email}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return m.User{}, err
	}

	user.Password = string(bytes)

	user, err = r.AddUser(user)

	if err != nil {
		return m.User{}, err
	}

	return user, nil
}

func GetUserFromToken(token string) m.User {
	claims := jwt.ParseToken(token)

	user := m.User{Email: claims.Email, Username: claims.Username, Id: claims.ID}

	return user
}

func GetUserFromGinContext()

func GetUserBooks(user m.User) []m.Book {
	books := r.GetUserBooks(user)

	return books;
}

func AddBookToUser(user m.User, bookId string) {
	r.AddBookToUser(user.Id, bookId)

	return
}
