package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/RichardHoa/blog-aggerator/internal/config"
	"github.com/RichardHoa/blog-aggerator/internal/database"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type AuthedHandler func(http.ResponseWriter, *http.Request, database.User)

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

func GetUserThroughAPIKey() AuthedHandler {
	return func(w http.ResponseWriter, r *http.Request, user database.User) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func CreateFeed(apiCfg *config.ApiConfig) AuthedHandler {

	return func(w http.ResponseWriter, r *http.Request, user database.User) {
		userID := user.ID

		// Get the request body
		var feedData map[string]string
		if err := json.NewDecoder(r.Body).Decode(&feedData); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid request body, it should be in json")
			return
		}

		// Extract the feed data
		name, ok := feedData["name"]
		if !ok || name == "" {
			RespondWithError(w, http.StatusBadRequest, "name field is required")
			return
		}

		url, ok := feedData["url"]
		if !ok || url == "" {
			RespondWithError(w, http.StatusBadRequest, "url field is required")
			return
		}

		ctx := r.Context()

		feed, err := apiCfg.DB.CreateFeed(ctx, database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
			Url:       url,
			UserID:    userID,
		})

		feedFollow, err := apiCfg.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FeedID:    feed.ID,
			UserID:    userID,
		})

		if err != nil {
			errString := err.Error()
			RespondWithError(w, http.StatusInternalServerError, "Failed to create feed: "+errString)
			return
		}

		var data = map[string]interface{}{
			"feed":        feed,
			"feed_follow": feedFollow,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}

}

func GetFeeds(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		feeds, err := apiCfg.DB.GetFeeds(ctx)
		if err != nil {
			errString := err.Error()
			RespondWithError(w, http.StatusInternalServerError, "Failed to get feeds: "+errString)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(feeds)
	}
}

func CreateFeedFollow(apiCfg *config.ApiConfig) AuthedHandler {
	return func(w http.ResponseWriter, r *http.Request, user database.User) {

		var feedData map[string]string
		if err := json.NewDecoder(r.Body).Decode(&feedData); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid request body, it should be in json")
			return
		}

		feedID, ok := feedData["feed_id"]
		if !ok || feedID == "" {
			RespondWithError(w, http.StatusBadRequest, "feed_id field is required")
			return
		}

		feedUUID, err := uuid.Parse(feedID)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid feed_id format")
			return
		}

		ctx := r.Context()
		feedFollow, err := apiCfg.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FeedID:    feedUUID,
			UserID:    user.ID,
		})
		if err != nil {
			errString := err.Error()
			RespondWithError(w, http.StatusInternalServerError, "Failed to create feed follow: "+errString)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(feedFollow)

	}
}

func DeleteFeedFollow(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the feedFollowID from the URL parameters
		urlPath := r.URL.Path
		parts := strings.Split(urlPath, "/")
		feedFollowID := parts[len(parts)-1]

		// vars := mux.Vars(r)
		// fmt.Printf("vars: %v\n", vars)
		// feedFollowID := vars["feedFollowID"]
		fmt.Printf("feedFollowID: %v\n", feedFollowID)

		feedFollowIDUUID, err := uuid.Parse(feedFollowID)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid feedFollowID format")
			return
		}

		deleteErr := apiCfg.DB.DeleteFeedFollow(r.Context(), feedFollowIDUUID)
		if deleteErr != nil {
			errString := deleteErr.Error()
			RespondWithError(w, http.StatusInternalServerError, "Failed to delete feed follow: "+errString)
			return
		}
		w.WriteHeader(http.StatusOK)

	}
}

func GetFeedFollows(apiCfg *config.ApiConfig) AuthedHandler {
	return func(w http.ResponseWriter, r *http.Request, user database.User) {
		ctx := r.Context()
		feedFollows, err := apiCfg.DB.GetFeedFollows(ctx, user.ID)
		if err != nil {
			errString := err.Error()
			RespondWithError(w, http.StatusInternalServerError, "Failed to get feed follows: "+errString)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(feedFollows)
	}
}

func GetPosts(apiCfg *config.ApiConfig) AuthedHandler {
	return func(w http.ResponseWriter, r *http.Request, user database.User) {
		ctx := r.Context()
		var allPosts []database.Post
		feeds, err := apiCfg.DB.GetFeedByUserID(ctx, user.ID)
		if err != nil {
			errString := err.Error()
			RespondWithError(w, http.StatusInternalServerError, "Failed to get feeds: "+errString)
			return
		}
		for _, feed := range feeds {

			posts, err := apiCfg.DB.GetPosts(ctx, database.GetPostsParams{
				FeedID: feed.ID,
				Limit:  10000})

			if err != nil {
				errString := err.Error()
				RespondWithError(w, http.StatusInternalServerError, "Failed to get posts: "+errString)
				return
			}
			allPosts = append(allPosts, posts...)
		}
		// Respond with all posts
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(allPosts); err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to encode posts: "+err.Error())
		}
	}
}
