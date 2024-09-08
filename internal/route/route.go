package route

import (
	"net/http"

	"github.com/RichardHoa/blog-aggerator/internal/config"
	"github.com/RichardHoa/blog-aggerator/internal/handlers"
	"github.com/RichardHoa/blog-aggerator/internal/middleware"
)

func ConfigureRoutes(mux *http.ServeMux, apiCfg *config.ApiConfig) {

	mux.HandleFunc("GET /v1/healhz", handlers.HandleHealthz)

	mux.HandleFunc("GET /v1/err", handlers.HandleReturnError)

	mux.HandleFunc("POST /v1/user", handlers.CreateUser(apiCfg))

	mux.HandleFunc("GET /v1/user", middleware.MiddlewareAuth(apiCfg, handlers.GetUserThroughAPIKey()))

	mux.HandleFunc("POST /v1/feed", middleware.MiddlewareAuth(apiCfg, handlers.CreateFeed(apiCfg)))

	mux.HandleFunc("GET /v1/feeds", handlers.GetFeeds(apiCfg))

	mux.HandleFunc("POST /v1/feed_follows",middleware.MiddlewareAuth(apiCfg, handlers.CreateFeedFollow(apiCfg)))

	mux.HandleFunc("DELETE /v1/feed_follows/", handlers.DeleteFeedFollow(apiCfg))

	mux.HandleFunc("GET /v1/feed_follows", middleware.MiddlewareAuth(apiCfg, handlers.GetFeedFollows(apiCfg)))
}
