package handlers

import (
	"encoding/json"
	"github.com/RichardHoa/blog-aggerator/internal/config"
	"github.com/RichardHoa/blog-aggerator/internal/database"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"time"
)

func HandleHealthz(w http.ResponseWriter, r *http.Request) {

	respondText := map[string]string{
		"status": "OK"}

	RespondWithJSON(w, http.StatusOK, respondText)

}

func HandleReturnError(w http.ResponseWriter, r *http.Request) {

	RespondWithError(w, 500, "Internal Server Error")
}

func CreateUser(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the request body
		var userData map[string]string
		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid request body, it should be in json")
			return
		}

		// Extract the user data
		name, ok := userData["name"]
		if !ok || name == "" {
			RespondWithError(w, http.StatusBadRequest, "name field is required")
			return
		}

		// Insert the user into the database
		ctx := r.Context()

		user, err := apiCfg.DB.CreateUser(ctx, database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
		})

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}

		// Return the created user
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

func GetUserThroughAPIKey(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondWithError(w, http.StatusBadRequest, "Authorization header is missing")
			return
		}

		// Assuming the format is "ApiKey <key>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "ApiKey" {
			RespondWithError(w, http.StatusBadRequest, "Invalid Authorization header format")
			return
		}

		apiKey := parts[1]

		ctx := r.Context()

		user, err := apiCfg.DB.GetUserByApiKey(ctx, apiKey)
		if err != nil {
			errStr := err.Error()
			RespondWithError(w, http.StatusInternalServerError, "Failed to get user: "+errStr)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}
