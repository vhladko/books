package models

import (
	"database/sql"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Email    string `json:"email"`
	Username sql.NullString `json:"username"`
	jwt.RegisteredClaims
}

type User struct {
	Id        string
	CreatedAt string
	Email     string
	Password  string
	Username  sql.NullString
}
