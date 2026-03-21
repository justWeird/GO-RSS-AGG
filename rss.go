package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

// we need a function to convert a url to feed. It accesses the website and reads the xml

// represent the type we're converting to based on existing rss xml structure
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RSSFeed, error) {
	// since we're accessing the url, we need a http client
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	// now access the url and get the entire response
	response, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}
	// close the response body after we're done with it to free up resources.
	// This is important to prevent memory leaks and ensure that the connection is properly closed.
	defer response.Body.Close()

	// now parse the entire body
	readData, err := io.ReadAll(response.Body)
	if err != nil {
		return RSSFeed{}, err
	}

	// now we have the data, we can unmarshal it into our RSSFeed struct
	rssFeed := RSSFeed{} //empty initialization
	err = xml.Unmarshal(readData, &rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}

	return rssFeed, nil

}
