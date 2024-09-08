package middleware

import (
		"net/http"
	"strings"
		"github.com/RichardHoa/blog-aggerator/internal/handlers"
		"github.com/RichardHoa/blog-aggerator/internal/config"
)

func MiddlewareAuth( cfg *config.ApiConfig,  handler handlers.AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusBadRequest)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "ApiKey" {
			http.Error(w, "Invalid Authorization header format", http.StatusBadRequest)
			return
		}

		apiKey := parts[1]

		ctx := r.Context()
		user, err := cfg.DB.GetUserByApiKey(ctx, apiKey)
		if err != nil {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		handler(w, r, user)
	}
}