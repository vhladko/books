package apihandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vhladko/books/helpers"
	"github.com/vhladko/books/repositories"
	"github.com/vhladko/books/services"
)

func HandleGetBookByIsbn(c *gin.Context) {
	isbn := c.Param("isbn")
	isbn, err := helpers.FormatIsbn(isbn)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err": err})
		return
	}
	book, err := services.GetBookByIsbn(isbn)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": book})
}

func HandleGetBookById(c *gin.Context) {
	id := c.Param("id")

	book, err := repositories.GetBookById(id)

	if err != nil {
		c.JSON(404, gin.H{"err": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"book": book})
}

func HandleGetBook(c *gin.Context) {
	isbn := c.Query("isbn")
	isbn, err := helpers.FormatIsbn(isbn)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err": err})
		return
	}
	book, err := services.GetBookByIsbn(isbn)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}
