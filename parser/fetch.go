package parser

import (
	"fmt"
	"v2/config"

	"github.com/mmcdole/gofeed"
)

func FetchAndStoreFeed(feedURL string, feedID int) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		fmt.Println("Error parsing feed:", err)
		return
	}

	for _, item := range feed.Items {
		var exists int
		err := config.DB.Get(&exists, "SELECT COUNT(*) FROM feed_items WHERE guid=$1", item.GUID)
		if err == nil && exists > 0 {
			continue
		}

		_, err = config.DB.Exec(`
			INSERT INTO feed_items (feed_id, title, link, published_at, guid, content)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, feedID, item.Title, item.Link, item.PublishedParsed, item.GUID, item.Content)

		if err != nil {
			fmt.Println("Insert error:", err)
		}
	}
	fmt.Println("update successful")
}
