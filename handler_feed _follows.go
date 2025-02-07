package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/JJFelix/rss_aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig)handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	type parameters struct{
		FeedId uuid.UUID `json:"feed_id"`
	}
	decode := json.NewDecoder(r.Body)

	params := parameters{}

	err := decode.Decode(&params)
	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: params.FeedId,
	})

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed follow: %v", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, databaseFeedFollowToFeedFollow(feed)) // empty struct will marshal into an empty JSON object
}

func (apiCfg *apiConfig)handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User){
	FeedFollow, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed follows: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedFollowsToFeedFollows(FeedFollow))
}

func (apiCfg *apiConfig)handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	feedFollwIDStr := chi.URLParam(r, "feedFollowID")
	FeedFollowID, err := uuid.Parse((feedFollwIDStr))
	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse feed follow id : %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{
		ID: FeedFollowID,
		UserID: user.ID,
	})
	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete feed follow id : %v", err))
		return
	}

	respondWithJSON(w, 200, struct{}{}) // return empty json object
}