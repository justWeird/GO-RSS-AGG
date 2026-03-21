package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/justWeird/GO-RSS-AGG/internal/database"
)

// This is a handler file for allowing a user to "subscribe"n to a feed i.e follow a feed.
// This creates an entry into the followed_feeds table, which has a user_id and a feed_id.
// This allows us to keep track of which feeds a user is following, and also allows us to
// easily query for all the feeds that a user is following when we want to generate their feed.

// since this route is auth protected, it will use middleware
func (db *dbConfig) handlerFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	// define the parameters to be decoded from the request bodt: feed id, user id (from auth middleware)
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err)) // using 400 (Bad Request) to indicate that the error is due to a client-side issue
		return
	}

	// call the database function to create a new followed feed entry
	followedFeed, err := db.DB.CreateFollowedFeed(r.Context(), database.CreateFollowedFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error following feed: %v", err))
		return
	}

	respondWithJSON(w, 200, dbFollowedFeedToFollowedFeed(followedFeed))

}

func (db *dbConfig) handlerGetFollowedFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	followedFeeds, err := db.DB.GetFollowedFeedsByUserID(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error fetching followed feeds: %v", err))
		return
	}

	respondWithJSON(w, 200, dbFollowedFeedsToFollowedFeeds(followedFeeds))

}

func (db *dbConfig) handlerUnfollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	// since we're passing a url parameter (instead of in the request body), we can extract it using chi's URLParam function.
	followedFeedIDStr := chi.URLParam(r, "followed_id")
	followedFeedID, err := uuid.Parse(followedFeedIDStr)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing followed feed ID: %v", err))
		return
	}

	err = db.DB.DeleteFollowedFeed(r.Context(), database.DeleteFollowedFeedParams{
		ID:     followedFeedID,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error unfollowing feed: %v", err))
		return
	}

	respondWithJSON(w, 201, map[string]string{"message": "Successfully unfollowed feed"})
}
