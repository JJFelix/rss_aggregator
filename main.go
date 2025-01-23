package main

import (
	// "fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main(){
	// Loading environment variables
	godotenv.Load(".env")

	portString := os.Getenv("PORT")

	if portString == ""{
		log.Fatal("PORT Not found")
	}
	

	// setting up a new router (Handler)
	router := chi.NewRouter()

	// cors configuration
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: 	[]string{"https://*", "http://*"}		,
		AllowedMethods: 	[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: 	[]string{"*"},
		ExposedHeaders: 		[]string{"*"},
		AllowCredentials: 	false,
		MaxAge: 			300,	
	}))

	// HTTP Server
	srv := &http.Server{
		Handler: router,
		Addr: ":"+ portString,
	}

	log.Printf("Server running on port %v", portString)

	err := srv.ListenAndServe() // a blocking operation
	if err != nil{
		log.Fatal(err)
	}
}