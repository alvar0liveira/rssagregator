package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alvar0liveira/rssagregator/internal/database"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
	}

	feed, err := apiConfig.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't Create Feed: %s", err))
		return
	}

	respondWithJson(w, 201, databaseFeedToFeed(feed))
}

func (apiConfig *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiConfig.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't Get Feeds: %s", err))
		return
	}

	respondWithJson(w, 201, databaseFeedsToFeeds(feeds))
}
