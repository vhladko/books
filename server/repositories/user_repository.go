package repositories

import (
	m "github.com/vhladko/books/models"
)

func AddUser(user m.User) (m.User, error) {
	insertUserQuery := `insert into "public"."user_"(email, password ) values($1, $2) returning id, created_at, email`
	err := DB.QueryRow(insertUserQuery, user.Email, user.Password).Scan(&user.Id, &user.CreatedAt, &user.Email)

	if err != nil {
		return user, err
	}

	return user, nil
}

func GetUserByEmail(email string) (m.User, error) {
	var user = m.User{}

	err := DB.
		QueryRow(`select id, created_at, email, username, password from "public"."user_" where "email"=$1`, email).
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
		QueryRow(`select id, created_at, email, username, password from "public"."user_" where "username"=$1`, username).
		Scan(&user.Id, &user.CreatedAt, &user.Email, &user.Username, &user.Password)

	if err != nil {
		return user, err
	} else {
		return user, nil
	}
}
