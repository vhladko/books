package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	r "github.com/vhladko/books/repositories"
	handlers "github.com/vhladko/books/apihandlers"
)

var DB *sql.DB

func main() {
	connectToDatabase()
	r := gin.Default()


	r.GET("/book/isbn/:isbn", handlers.HandleAuthGuard, handlers.HandleGetBookByIsbn)
	r.GET("/book/id/:id", handlers.HandleAuthGuard, handlers.HandleGetBookById)
	r.POST("/login", handlers.HandleLogin)
	r.GET("/logout", handlers.HandleLogout)
	r.POST("/register", handlers.HandleRegister)
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
