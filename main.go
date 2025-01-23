package main

import (
	// "fmt"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JJFelix/rss_aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct{
	DB *database.Queries
}

func main(){
	// Loading environment variables
	godotenv.Load(".env")

	portString := os.Getenv("PORT")

	if portString == ""{
		log.Fatal("PORT Not found")
	}

	// importimg the db connection
	dbURL := os.Getenv("DB_URL")

	if dbURL == ""{
		log.Fatal("DB_URL Not found")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil{
		log.Fatal("Can't connect to Database:", err)
	}

	db := database.New(conn)
	
	apiCfg := apiConfig{
		DB: db,
	}

	// Start of Scrapping
	go startScraping(db, 10, time.Minute)	


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

	// routers
	v1Router := chi.NewRouter() // create new router instance for version 1 of the app

	v1Router.Get("/ready", handlerReadiness) // connecting the handlerReadiness function to the ready path(endpoint)
	v1Router.Get("/err", handlerErr)

	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	// mount the v1 router to the main router under v1 prefix
	// routes defined under v1Router will be accessible under the '/v1' prefix
	// e.g. /v1/ready
	// allows for scalability, modularity and API Versioning
	router.Mount("/v1", v1Router)

	// HTTP Server
	srv := &http.Server{
		Handler: router,
		Addr: ":"+ portString,
	}

	log.Printf("Server running on port %v", portString)

	err = srv.ListenAndServe() // a blocking operation
	if err != nil{
		log.Fatal(err)
	}
}