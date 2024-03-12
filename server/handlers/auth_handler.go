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
		c.JSON(http.StatusNotFound, gin.H{"err": "User not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Wrong credentials"})
		return
	}

	expirationTime := jwt.GetExparationTime()
	token := jwt.CreateToken(user, expirationTime)

	signedString, err := jwt.SignToken(token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Wasnt able to generate token"})
	}

	c.SetCookie("token", signedString, int(expirationTime.Unix()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, "")
}

func HandleLogout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, "")
}

func HandleRegister(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	_, err := r.GetUserByEmail(email)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "User with such email already exists"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Was not able to connect to database"})
		return
	}

	user, err := services.AddUser(email, password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Was not able to create a user"})
		return
	}

	expirationTime := jwt.GetExparationTime()
	token := jwt.CreateToken(user, expirationTime)

	signedString, err := jwt.SignToken(token)

	if err != nil {
		fmt.Print(err, "err")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Wasn't able to generate token"})
		return
	}

	c.SetCookie("token", signedString, int(expirationTime.Unix()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, "")
}

func HandleAuthGuard(c *gin.Context) {
	token, err := c.Cookie("token")

	isEmpty := token == ""
	fmt.Printf("token %s and its %t", token, isEmpty)
	
	if err != nil || token == "" {
		fmt.Print("ERROR HERE")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Wasnt able to find token"})
		return
	}

	fmt.Print("FOR SOME REASON IM HERE")

	err = jwt.VerifyToken(token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Invalid token"})
		return
	}

	user := services.GetUserFromToken(token);

	c.Set("user", user)

	c.Next()
}
