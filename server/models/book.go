package models

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
