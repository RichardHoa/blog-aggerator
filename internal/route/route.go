package route

import (
	"net/http"

	"github.com/RichardHoa/blog-aggerator/internal/handlers"
	"github.com/RichardHoa/blog-aggerator/internal/config"
)

func ConfigureRoutes(mux *http.ServeMux, apiCfg *config.ApiConfig) {

	mux.HandleFunc("GET /v1/healhz",handlers.HandleHealthz )

	mux.HandleFunc("GET /v1/err", handlers.HandleReturnError)

	mux.HandleFunc("POST /v1/user", handlers.CreateUser(apiCfg))

	mux.HandleFunc("GET /v1/user", handlers.GetUserThroughAPIKey(apiCfg))
}