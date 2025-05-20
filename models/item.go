package models

import "time"

type FeedItem struct {
	ID        int       `db:"id" json:"id"`
	FeedID    int       `db:"feed_id" json:"feed_id"`
	Title     string    `db:"title" json:"title"`
	Link      string    `db:"link" json:"link"`
	Published time.Time `db:"published_at" json:"published_at"`
	GUID      string    `db:"guid" json:"guid"`
	Content   string    `db:"content" json:"content"`
}
