package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/justWeird/GO-RSS-AGG/internal/database"
)

/*
a backgraound scraper functoin that runs in the background and periodically
fetches the latest entries from the feeds that users are following, and updates
the database with the new entries. This function will run in a separate goroutine
when the application starts, and will continue to run indefinitely,
fetching new entries every hour or so.
*/

// Because it's a background task, it will run in a separate goroutine,
// which allows it to run concurrently with the main server without
// blocking incoming HTTP requests.
func startBackgroundScraper(
	db *database.Queries, // pass in the database connection so that the scraper can query for the feeds that users are following and update the database with new entries
	concurrency int, // specify the number of concurrent workers to use for fetching feeds. This allows us to fetch multiple feeds in parallel, improving the efficiency of the scraper.
	timeBetweenRequest time.Duration, // specify the time to wait between requests to the same feed
) {
	// since it's gonna run indefinitely, we can use a for loop to keep it running
	// important to have a logger so we know what's happening in the scraper
	log.Printf("Scraper started with concurrency: %d and time between requests: %s", concurrency, timeBetweenRequest)

	// use a ticker to run the scraper at regular intervals
	ticker := time.NewTicker(timeBetweenRequest) // run the scraper every hour

	// the ticker uses a channel to send a signal every time the specified duration has passed.
	//  We can use this channel to trigger the scraper function to run at regular intervals.
	for ; ; <-ticker.C { //this syntax allows the ticker to trigger immediately on the first run.
		log.Println("Scraper tick: fetching feeds...")
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Error fetching feeds: %v", err)
			continue //continue not return because the scraper runs indefinitely.
		}

		// we need to fetch the entries in each each feed concurrently.
		// use a wait group to wait before moving onto the next part. Gotten from the STL
		waitGroup := &sync.WaitGroup{}

		for _, feed := range feeds {
			waitGroup.Add(1) // increment the wait group counter for each feed

			// fetch the feed in a separate goroutine to allow for concurrent fetching of multiple feeds.
			go scrapeFeed(waitGroup, db, feed)
		}

		waitGroup.Wait() // wait for all feed fetching goroutines to finish before moving on to the next tick of the scraper

	}

}

func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done() // decrement the wait group counter when the function returns

	// here we would implement the logic to fetch the feed, parse the entries, and update the database with any new entries.
	err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed: %v", err)
		return
	}

	// for now, log all titles
	for _, entry := range rssFeed.Channel.Item {

		// becuase the description may or may not be empty, the use Nullstring to handle the case where the description is empty.
		// This allows us to insert a NULL value into the database when the description is empty, instead of an empty string.
		description := sql.NullString{}

		if entry.Description != "" {
			description.String = entry.Description
			description.Valid = true
		}

		// logic for parsing the publication date from the RSS feed, which is typically in a string format.
		pubDate, err := time.Parse(time.RFC1123Z, entry.PubDate)
		if err != nil {
			log.Printf("Error parsing publication date: %v", err)
			return
		}

		// insert into the database
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       entry.Title,
			Url:         entry.Link,
			Description: description,
			PublishedAt: pubDate, // this is a string in the RSS feed, but we can parse it into a time.Time object before inserting it into the database.
			FeedID:      feed.ID,
		})
		if err != nil {
			// a possible error is that the post already exists in the database, which can happen if the scraper runs multiple times and encounters the same entries.
			//  In this case, we can log the error and continue with the next entry, instead of returning and stopping the entire scraping process.
			if strings.Contains(err.Error(), "duplicate key value") {
				continue
			}
			log.Printf("Error creating post: %v", err)
			continue
		}
	}

	log.Printf("Fetched %d entries from feed: %s", len(rssFeed.Channel.Item), feed.Name)

}
