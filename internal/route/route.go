package route

import (
	"net/http"

	"github.com/RichardHoa/blog-aggerator/internal/handlers"
)

func ConfigureRoutes(mux *http.ServeMux) {

	mux.HandleFunc("GET /v1/healhz",handlers.HandleHealthz )

	mux.HandleFunc("GET /v1/err", handlers.HandleReturnError)
}