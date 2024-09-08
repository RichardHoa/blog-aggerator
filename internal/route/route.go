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
}
