// package handlers

// import (
// 	"net/http"
// 	"v2/config"
// 	"v2/models"
// 	"v2/parser"

// 	"github.com/gin-gonic/gin"
// )

// func AddFeed(c *gin.Context) {
// 	var feed models.Feed
// 	if err := c.ShouldBindJSON(&feed); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err := config.DB.QueryRow(
// 		"INSERT INTO feeds (url, title, description) VALUES ($1, $2, $3) RETURNING id",
// 		feed.URL, feed.Title, feed.Description).Scan(&feed.ID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add feed"})
// 		return
// 	}

// 	go parser.FetchAndStoreFeed(feed.URL, feed.ID)

// 	c.JSON(http.StatusOK, feed)
// }

// func GetFeeds(c *gin.Context) {
// 	var feeds []models.Feed
// 	err := config.DB.Select(&feeds, "SELECT * FROM feeds")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch feeds"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, feeds)
// }

// func GetFeedItems(c *gin.Context) {
// 	feedID := c.Param("id")
// 	var items []models.FeedItem
// 	err := config.DB.Select(&items, "SELECT * FROM feed_items WHERE feed_id = $1 ORDER BY published_at DESC", feedID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, items)
// }

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"v2/config"
	"v2/middleware"
	"v2/models"
	"v2/parser"
)

// AddFeed handles POST /feeds
// func AddFeed(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var feed models.Feed
// 	if err := json.NewDecoder(r.Body).Decode(&feed); err != nil {
// 		http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 		return
// 	}

// 	userID := middleware.GetUserID(r)
// 	feed.UserID = userID

// 	err := config.DB.QueryRow(
// 		"INSERT INTO feeds (url, title, description, user_id) VALUES ($1, $2, $3, $4) RETURNING id",
// 		feed.URL, feed.Title, feed.Description, feed.UserID).Scan(&feed.ID)
// 	if err != nil {
// 		http.Error(w, "Failed to add feed", http.StatusInternalServerError)
// 		return
// 	}

// 	go parser.FetchAndStoreFeed(feed.URL, feed.ID)

// 	writeJSON(w, feed)
// }

func AddFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var feed models.Feed
	if err := json.NewDecoder(r.Body).Decode(&feed); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r)
	feed.UserID = userID

	// Try to insert the feed
	err := config.DB.QueryRow(
		"INSERT INTO feeds (url, title, description, user_id) VALUES ($1, $2, $3, $4) RETURNING id",
		feed.URL, feed.Title, feed.Description, feed.UserID).Scan(&feed.ID)

	if err != nil {
		// If insertion failed, check if the feed already exists for the user
		// query := "SELECT id FROM feeds WHERE url = $1 AND user_id = $2"
		query := "SELECT id FROM feeds WHERE url = $1"
		err2 := config.DB.QueryRow(query, feed.URL).Scan(&feed.ID)
		if err2 != nil {
			log.Println(err2)
			http.Error(w, "Failed to add or retrieve feed", http.StatusInternalServerError)
			return
		}
	}

	// Fetch and store items ONLY if:
	//  - Insert was successful, OR
	//  - Feed already existed (retrieved above)
	go parser.FetchAndStoreFeed(feed.URL, feed.ID)

	writeJSON(w, feed)
}

// GetFeeds handles GET /feeds
func GetFeeds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)

	var feeds []models.Feed
	err := config.DB.Select(&feeds, "SELECT * FROM feeds WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Failed to fetch feeds", http.StatusInternalServerError)
		return
	}

	writeJSON(w, feeds)
}

// GetFeedItems handles GET /feeds/{id}/items
func GetFeedItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract feed ID from URL manually
	path := strings.TrimPrefix(r.URL.Path, "/feeds/")
	idPart := strings.Split(path, "/")[0]
	feedID, err := strconv.Atoi(idPart)
	if err != nil {
		http.Error(w, "Invalid feed ID", http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r)

	// Ensure feed belongs to the logged-in user
	var ownerID int
	err = config.DB.Get(&ownerID, "SELECT user_id FROM feeds WHERE id = $1", feedID)
	if err != nil || ownerID != userID {
		http.Error(w, "Not authorized or feed not found", http.StatusForbidden)
		return
	}

	var items []models.FeedItem
	err = config.DB.Select(&items, "SELECT * FROM feed_items WHERE feed_id = $1 ORDER BY published_at DESC", feedID)
	if err != nil {
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
		return
	}

	writeJSON(w, items)
}

// writeJSON sends a JSON response
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
