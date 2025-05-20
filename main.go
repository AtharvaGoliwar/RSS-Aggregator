// package main

// import (
// 	"os"
// 	"v2/config"
// 	"v2/routes"
// 	"v2/scheduler"
// )

// func main() {
// 	config.InitDB()
// 	scheduler.StartFeedUpdater()
// 	r := routes.SetupRouter()
// 	r.Run(":" + os.Getenv("PORT"))
// }

// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"strings"
// 	"v2/config"
// 	"v2/handlers"
// 	"v2/middleware"
// 	"v2/scheduler"
// )

// func main() {
// 	config.InitDB()
// 	scheduler.StartFeedUpdater()

// 	http.HandleFunc("/feeds", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodPost {
// 			handlers.AddFeed(w, r)
// 		} else if r.Method == http.MethodGet {
// 			handlers.GetFeeds(w, r)
// 		} else {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		}
// 	})

// 	http.HandleFunc("/signup", handlers.Signup)
// 	http.HandleFunc("/login", handlers.Login)

// 	http.HandleFunc("/feeds", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
// 		switch r.Method {
// 		case http.MethodPost:
// 			handlers.AddFeed(w, r)
// 		case http.MethodGet:
// 			handlers.GetFeeds(w, r)
// 		default:
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		}
// 	}))

// 	http.HandleFunc("/feeds/", middleware.Auth(handlers.GetFeedItems))

// 	http.HandleFunc("/feeds/", func(w http.ResponseWriter, r *http.Request) {
// 		if strings.HasSuffix(r.URL.Path, "/items") {
// 			handlers.GetFeedItems(w, r)
// 		} else {
// 			http.Error(w, "Not found", http.StatusNotFound)
// 		}
// 	})

// 	fmt.Println("Server running on http://localhost:8080")
// 	http.ListenAndServe(":8080", nil)
// }

package main

import (
	"log"
	"net/http"
	"v2/config"
	"v2/handlers"
	"v2/middleware"

	_ "github.com/lib/pq"
)

func main() {
	config.InitDB()

	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)

	http.Handle("/feeds", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetFeeds)))
	http.Handle("/feeds/", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetFeedItems)))
	http.Handle("/add-feed", middleware.AuthMiddleware(http.HandlerFunc(handlers.AddFeed)))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
