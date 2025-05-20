package models

type Feed struct {
	ID          int    `db:"id" json:"id"`
	URL         string `db:"url" json:"url"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	UserID      int    `json:"user_id" db:"user_id"`
}
