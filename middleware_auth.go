package main

import (
	"fmt"
	"net/http"

	"github.com/JJFelix/rss_aggregator/internal/database"
	"github.com/JJFelix/rss_aggregator/internal/database/auth"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil{
		respondWithError(w, http.StatusForbidden, fmt.Sprintf("Auth Error: %v", err))
		return 
	}
	user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
		return 
	}
	handler(w, r, user)
	}
}