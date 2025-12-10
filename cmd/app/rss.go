package main

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (a *Application) FetchFeedHandler(cmd Command) error {
	if len(cmd.Args) < 1 {
		return ErrNoArgumentsProvided
	}

	if len(cmd.Args) > 1 {
		return ErrTooManyArguments
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	feedURL := cmd.Args[0]

	rssFeed, err := fetchFeed(ctx, feedURL)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("%+v", rssFeed)
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	// Identify program to server
	req.Header.Add("User-Agent", "gator")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rf := &RSSFeed{}
	err = xml.Unmarshal(body, rf)
	if err != nil {
		return nil, err
	}

	return rf, nil
}
