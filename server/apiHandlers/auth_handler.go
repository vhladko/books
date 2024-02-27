package apihandlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vhladko/books/jwt"
	r "github.com/vhladko/books/repositories"
	"github.com/vhladko/books/services"
	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := r.GetUserByEmail(email)

	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusNotFound, gin.H{"err": "user not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "wrong credentials"})
		return
	}

	expirationTime := jwt.GetExparationTime()
	token := jwt.CreateToken(user, expirationTime)

	signedString, err := jwt.SignToken(token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "wasnt able to generate token"})
	}

	c.SetCookie("token", signedString, int(expirationTime.Unix()), "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func HandleLogout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

func HandleRegister(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	_, err := r.GetUserByEmail(email)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user with such email already exists"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "was not able to connect to database", "exact": err})
		return
	}

	user, err := services.AddUser(email, password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "was not able to create a user"})
		return
	}

	expirationTime := jwt.GetExparationTime()
	token := jwt.CreateToken(user, expirationTime)

	signedString, err := jwt.SignToken(token)

	if err != nil {
		fmt.Print(err, "err")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "wasnt able to generate token", "exact": err})
		return
	}

	c.SetCookie("token", signedString, int(expirationTime.Unix()), "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func HandleAuthGuard(c *gin.Context) {
	token, err := c.Cookie("token")

	if err != nil || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "wasnt able to find token"})
		return
	}

	err = jwt.VerifyToken(token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid token"})
		return
	}

	c.Next()
}
