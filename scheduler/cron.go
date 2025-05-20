package scheduler

import (
	"fmt"
	"v2/config"
	"v2/parser"

	"github.com/robfig/cron/v3"
)

func StartFeedUpdater() {
	c := cron.New()
	c.AddFunc("@every 10m", func() {
		var feeds []struct {
			ID  int    `db:"id"`
			URL string `db:"url"`
		}
		err := config.DB.Select(&feeds, "SELECT id, url FROM feeds")
		if err != nil {
			fmt.Println("Cron fetch error:", err)
			return
		}

		for _, feed := range feeds {
			go parser.FetchAndStoreFeed(feed.URL, feed.ID)
		}
	})
	c.Start()
}
