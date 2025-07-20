package entity

import "time"

type Advertise struct {
	Id          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Image       string    `db:"image" json:"image"`
	Price       float64   `db:"price" json:"price"`
	UserId      int       `db:"user_id" json:"user_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
