package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	_ "github.com/lib/pq"
)

const goodreadsUrl = "https://goodreads.com/book/isbn/"

type Book struct {
	Id          string `json:"id"`
	CreatedAt   string `json:"createdAt"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TotalPages  int16  `json:"totalPages"`
	Isbn        string `json:"isbn"`
	PreviewUrl  string `json:"previewUrl"`
	AuthorId    string `json:"authorId"`
	AuthorName  string `json:"authorName"`
}

type Author struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	Name      string `json:"name"`
}

var db *sql.DB

func main() {
	connectToDatabase()
	r := gin.Default()

	r.GET("/book/isbn/:isbn", handleGetBookByIsbn)
	r.GET("/book/id/:id", handleGetBookById)
	r.Run()
}

func handleGetBookByIsbn(c *gin.Context) {
	isbn := c.Param("isbn")
	isbn, err := formatIsbn(isbn)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err": err})
		return
	}
	book, err := getBookByIsbn(isbn)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func handleGetBookById(c *gin.Context) {
	id := c.Param("id")

	book, err := getBookById(id)

	if err != nil {
		c.JSON(404, gin.H{"err": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func getBookById(id string) (Book, error) {
	book, err := findBookInDatabaseById(id)

	if err != nil {
		return Book{}, err
	}

	return book, err
}

func getBookByIsbn(isbn string) (Book, error) {
	book, err := findBookInDatabaseByIsbn(isbn)
	fmt.Println(book, err, "here")
	if err == nil {
		return book, nil
	}

	book, err = findBookOnGoodreads(isbn)

	if err != nil {
		return Book{}, err
	}

	author, err := getAuthorForGoodreads(book)
	book.AuthorId = author.Id
	instertBookIntoDatabase(&book)

	return book, err
}

func getAuthorByName(name string) (Author, error) {
	var author = Author{}
	err := db.
		QueryRow(`select id from "public"."author" where "name"=$1`, name).
		Scan(&author.Id)

	if err != nil {
		return author, err
	} else {
		return author, nil
	}
}

func instertAuthorIntoDatabase(name string) (Author, error) {
	var author = Author{}
	insertBookQuery := `insert into "public"."author"("name") values($1) returning id, created_at, name`
	err := db.QueryRow(insertBookQuery, name).Scan(&author.Id, &author.CreatedAt, &author.Name)

	if err != nil {
		return author, err
	}

	return author, nil

}

func getAuthorForGoodreads(book Book) (Author, error) {
	author, err := getAuthorByName(book.AuthorName)

	if err == nil {
		return author, err
	}

	author, err = instertAuthorIntoDatabase(book.AuthorName)

	return author, err

}

func connectToDatabase() {
	var err error
	urlExample := "postgres://vladhladko@localhost:5432/books?sslmode=disable"
	db, err = sql.Open("postgres", urlExample)
	if err != nil {
		panic(err)
	}

}

func findBookInDatabaseByIsbn(isbn string) (Book, error) {
	var book = Book{Isbn: isbn}

	err := db.
		QueryRow(`select id, created_at, title, author_id, author_name, total_pages, description, preview_url, isbn from "public"."book" where "isbn"=$1`, isbn).
		Scan(&book.Id, &book.CreatedAt, &book.Title, &book.AuthorId, &book.AuthorName, &book.TotalPages, &book.Description, &book.PreviewUrl, &book.Isbn)

	fmt.Println(book, "from db")

	if err != nil {
		return book, err
	} else {
		return book, nil
	}
}

func findBookInDatabaseById(id string) (Book, error) {
	var book = Book{Id: id}

	err := db.
		QueryRow(`select id, created_at, title, author_id, author_name, total_pages, description, preview_url, isbn from "public"."book" where "id"=$1`, id).
		Scan(&book.Id, &book.CreatedAt, &book.Title, &book.AuthorId, &book.AuthorName, &book.TotalPages, &book.Description, &book.PreviewUrl, &book.Isbn)

	if err != nil {
		return book, err
	} else {
		return book, nil
	}
}

func instertBookIntoDatabase(book *Book) {
	insertBookQuery := `insert into "public"."book"("title", "total_pages", "description", "preview_url", "isbn", "author_id", "author_name") values($1, $2, $3, $4, $5, $6, $7) returning id, created_at`
	e := db.QueryRow(insertBookQuery, book.Title, book.TotalPages, book.Description, book.PreviewUrl, book.Isbn, book.AuthorId, book.AuthorName).Scan(&book.Id, &book.CreatedAt)

	if e != nil {
		panic(e)
	}

}

func findBookOnGoodreads(isbn string) (Book, error) {
	var url = goodreadsUrl + isbn
	var book = Book{Isbn: isbn}
	var err error

	c := colly.NewCollector()

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println(r, e)
		err = e
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL.String())
	})

	c.OnHTML(".BookCover__image", func(e *colly.HTMLElement) {
		book.PreviewUrl = e.ChildAttr("img", "src")
	})

	c.OnHTML(".BookPageMetadataSection", func(e *colly.HTMLElement) {
		book.Description = e.ChildText("[data-testid='description']")
		pagesFormat := e.ChildText("[data-testid='pagesFormat']")
		pagesString := strings.Map(filterOnlyDigits, pagesFormat)
		totalPages, err := strconv.ParseInt(pagesString, 10, 16)
		if err != nil {
			book.TotalPages = 1
		} else {
			book.TotalPages = int16(totalPages)
		}

	})

	c.OnHTML(".ContributorLink__name", func(e *colly.HTMLElement) {
		book.AuthorName = e.Text
	})

	c.OnHTML("[data-testid='bookTitle']", func(e *colly.HTMLElement) {
		book.Title = e.Text
	})

	c.Visit(url)

	return book, err
}

func formatIsbn(isbn string) (string, error) {
	str := strings.Map(filterOnlyDigits, isbn)
	strLen := len(str)
	if strLen != 10 && strLen != 13 {
		return "", errors.New("invalid isbn")
	}
	return str, nil
}

func filterOnlyDigits(r rune) rune {
	if unicode.IsDigit(r) {
		return r
	} else {
		return -1
	}
}

type User struct {
	Id        string
	CreatedAt string
	Email     string
	Password  string
	Username  string
}

func createUser(email string, username string, password string) (User, error) {
	user := User{Email: email, Username: username}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	insertUserQuery := `insert into "public"."author"(email, username, password ) values($1, $2, $3) returning id, created_at, email, username`
	err = db.QueryRow(insertUserQuery, user.Email, user.Username, bytes).Scan(&user.Id, &user.CreatedAt, &user.Email, &user.Username)

	if err != nil {
		return user, err
	}

	return user, nil
}

func getUserByEmail(email string) (User, error) {

	var user = User{}

	err := db.
		QueryRow(`select id, created_at, email, username, password from "public"."user" where "email"=$1`, email).
		Scan(&user.Id, &user.CreatedAt, &user.Email, &user.Username, &user.Password)

	if err != nil {
		return user, err
	} else {
		return user, nil
	}
}

type MyClaim struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func handleLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := getUserByEmail(email)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "user not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "wrong credentials"})
		return
	}


	expirationTime := jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	token := createToken(user, expirationTime)

	signedString, err := signToken(token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "wasnt able to generate token"})
		return
	}

	c.SetCookie("token", signedString, int(expirationTime.Unix()), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"msg": "login successful"})
}

var secretKey = "secret-key"

func createToken(user User, expirationTime *jwt.NumericDate) *jwt.Token {
	claims := MyClaim{
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token;
}

func signToken(token *jwt.Token) (string, error) {
	signedString, err := token.SignedString(secretKey)

	return signedString, err
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil

}

func handleLogout(c *gin.Context) {

}

func handleAuthGuard(c *gin.Context) {
	token, err := c.Cookie("token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "wasnt able to find token"})
	}


	fmt.Print(token)


}
