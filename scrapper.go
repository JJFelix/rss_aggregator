package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/JJFelix/rss_aggregator/internal/database"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration){
	log.Printf("Scrapping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C{ // using channels
		feeds, err := db.GetNetxFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil{
			log.Println("error fetching feeds:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds{
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait() // blocking ooperation
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed){
	defer wg.Done() // to be called at the end of this function

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil{
		log.Println("error marking feeds as fetched feeds:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil{
		log.Println("Error fetching feed: ",err)
		return
	}

	for _, item := range rssFeed.Channel.Item{
		// log.Println("Found Post:", item.Title, "on feed", feed.Name)
		description := sql.NullString{}
		if item.Description != ""{
			description.String = item.Description
			description.Valid = true
		}
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil{
			log.Printf("Couldn't parse date %v with err %v", item.PubDate,err)
			continue
		}

		_, err = db.CreatePost(
			context.Background(), 		
			database.CreatePostParams{
				ID: 		uuid.New(),
				CreatedAt: 	time.Now().UTC(),
				UpdatedAt: 	time.Now().UTC(),
				Title: 		item.Title,
				Description: description,
				PublishedAt: pubAt,
				Url: item.Link,
				FeedID: feed.ID,
		})
		if err != nil{
			if strings.Contains(err.Error(), "duplicate key"){
				continue
			}
			log.Printf("Couldn't create post with err %v", err)
			// continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}