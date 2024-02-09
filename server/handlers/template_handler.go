package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vhladko/books/services"
	"github.com/vhladko/books/templates"
)


func HandleMain(c *gin.Context) {
	token, err := c.Cookie("token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized"})
		return
	}

	user := services.GetUserFromToken(token);
	c.HTML(http.StatusOK,"", templates.Home(user))
}

func HandleGetLogin(c *gin.Context) {
	c.HTML(http.StatusOK,"", templates.Login())
}

func HandleGetRegister(c *gin.Context) {
	c.HTML(http.StatusOK,"", templates.Register())
}
