package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/vhladko/books/handlers"
	templRender "github.com/vhladko/books/render"
	r "github.com/vhladko/books/repositories"
)

var DB *sql.DB

func main() {
	connectToDatabase()
	r := gin.Default()

	r.HTMLRender = templRender.Default;

	r.Static("/assets", "./assets")

	r.GET("/", handlers.HandleMain)
	r.GET("/book", handlers.HandleAuthGuard, handlers.HandleGetBook)
	r.GET("/book/isbn/:isbn", handlers.HandleAuthGuard, handlers.HandleGetBookByIsbn)
	r.GET("/book/id/:id", handlers.HandleAuthGuard, handlers.HandleGetBookById)
	r.GET("/login", handlers.HandleGetLogin)
	r.POST("/login", handlers.HandleLogin)
	r.GET("/logout", handlers.HandleLogout)
	r.POST("/register", handlers.HandleRegister)
	r.GET("/register", handlers.HandleGetRegister)
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
