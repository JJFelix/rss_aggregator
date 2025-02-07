package main

import (
	"encoding/json"
	"log"
	"net/http"
)
// utility error function to write the error messages for HTTP requests in JSON formart
func respondWithError(w http.ResponseWriter, code int, msg string){
	if code > 499{
		log.Println("Responding with 5XX eror:", msg)
	}

	type errorResponse struct{
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}


// utility function to write the resonse body for HTTP requests
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	data, err := json.Marshal(payload) // converts the payload into JSON
	if err != nil{
		log.Printf("Failed to marshal JSON Response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}