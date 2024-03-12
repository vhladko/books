package apihandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vhladko/books/models"
	"github.com/vhladko/books/services"
)

func HandleGetUserBooks(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	books := services.GetUserBooks(user)

	c.JSON(http.StatusOK, gin.H{"books": books})
}

func HandleAddBookToUser(c *gin.Context) {
	bookId := c.PostForm("bookId")
	user := c.MustGet("user").(models.User)

	services.AddBookToUser(user, bookId)

	c.JSON(http.StatusOK, "")
}
