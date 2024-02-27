package services

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/vhladko/books/helpers"
	m "github.com/vhladko/books/models"
)

const goodreadsUrl = "https://goodreads.com/book/isbn/"

func GetBookFromGoodreads(isbn string) (m.Book, error) {
	var url = goodreadsUrl + isbn
	var book = m.Book{Isbn: isbn}
	var err error

	c := colly.NewCollector()

	c.OnError(func(r *colly.Response, e error) {
		err = e
	})

	c.OnHTML(".BookCover__image", func(e *colly.HTMLElement) {
		book.PreviewUrl = e.ChildAttr("img", "src")
	})

	c.OnHTML(".BookPageMetadataSection", func(e *colly.HTMLElement) {
		pagesFormat := e.ChildText("[data-testid='pagesFormat']")
		pagesString := strings.Map(helpers.FilterOnlyDigits, pagesFormat)
		totalPages, err := strconv.ParseInt(pagesString, 10, 16)
		if err != nil {
			book.TotalPages = 1
		} else {
			book.TotalPages = int16(totalPages)
		}
	})

	c.OnHTML("[data-testid='description'] .Formatted", func(e *colly.HTMLElement) {
		desc, err := e.DOM.Html()
		if(err != nil) {
			book.Description = e.Text
		}
		book.Description = desc
	})

	c.OnHTML(".AuthorPreview .ContributorLink__name", func(e *colly.HTMLElement) {
		book.AuthorName = e.Text
	})

	c.OnHTML("[data-testid='bookTitle']", func(e *colly.HTMLElement) {
		book.Title = e.Text
	})

	c.Visit(url)

	return book, err
}
