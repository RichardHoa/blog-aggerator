package main

import (
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"context"
	"github.com/RichardHoa/blog-aggerator/internal/config"
	"github.com/RichardHoa/blog-aggerator/internal/database"
)

// RSSFeed represents the structure of an RSS feed
type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Items       []Item `xml:"item"`
	} `xml:"channel"`
}

// Item represents an individual item in the RSS feed
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchDataFromFeed(url string) (*RSSFeed, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}

func feedWorker(stop chan struct{}, fetchInterval time.Duration, numFeeds int32, apiConfig *config.ApiConfig) {
	log.Println("Feed worker started")

	ticker := time.NewTicker(fetchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			log.Println("Feed worker stopped")
			return
		case <-ticker.C:
			log.Println("Fetching feeds...")
			ctx := context.Background()

			feeds, err := apiConfig.DB.GetNextFeedsToFetch(ctx, numFeeds)
			if err != nil {
				log.Println("Error getting feeds:", err)
				continue
			}

			var wg sync.WaitGroup
			for _, feed := range feeds {
				wg.Add(1)
				go func(feed database.Feed) {
					defer wg.Done()
					log.Printf("Fetching feed %s...\n", feed.Url)

					RssFeed, err := FetchDataFromFeed(feed.Url)
					if err != nil {
						log.Printf("Error fetching feed %s: %v\n", feed.Url, err)
						return
					}

					log.Printf("Processing feed %s...\n", feed.Url)

					log.Printf("Feed title: %s\n", RssFeed.Channel.Title)

					lastFetched := sql.NullTime{
						Time:  time.Now(),
						Valid: true,
					}

					apiConfig.DB.MarkFeedFetched(
						ctx,
						database.MarkFeedFetchedParams{
							ID:          feed.ID,
							UpdatedAt:   time.Now(),
							LastFetched: lastFetched,
						})

					log.Println("Feed marked as fetched")

					// Items := RssFeed.Channel.Items
					// for _, post := range Items {
					// 	log.Println(post.Title)
					// }
				}(feed)
			}

			wg.Wait()
			log.Println("Feeds processed")
		}
	}
}
