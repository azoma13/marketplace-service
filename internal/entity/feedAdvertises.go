package entity

type FeedAdvertises struct {
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	ImageUrl       string  `json:"image_url"`
	Price          float64 `json:"price"`
	AuthorUsername string  `json:"author_username"`
	IsAuthor       bool    `json:"is_author,omitempty"`
}
