package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	r "github.com/vhladko/books/repositories"
	handlers "github.com/vhladko/books/handlers"
)

var DB *sql.DB

func main() {
	connectToDatabase()
	r := gin.Default()


	r.GET("/api/book/isbn/:isbn", handlers.HandleAuthGuard, handlers.HandleGetBookByIsbn)
	r.GET("/api/book/id/:id", handlers.HandleAuthGuard, handlers.HandleGetBookById)
	r.POST("/api/login", handlers.HandleLogin)
	r.GET("/api/logout", handlers.HandleLogout)
	r.POST("/api/register", handlers.HandleRegister)
	r.GET("/api/books", handlers.HandleAuthGuard, handlers.HandleGetUserBooks)
	r.POST("/api/book", handlers.HandleAuthGuard, handlers.HandleAddBookToUser)
	r.Run()
}

func connectToDatabase() {
	var err error
	url := "postgres://vladhladko@localhost:5432/books?sslmode=disable"
	r.DB, err = sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
}
