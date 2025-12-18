package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"jaytaylor.com/html2text"
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

	t, err := time.ParseDuration(cmd.Args[0])
	if err != nil{
		log.Fatal(err)
	}

	ticker := time.NewTicker(t)
	for ; ; <-ticker.C{
		a.scrapeFeeds()
	}

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

func (a *Application) scrapeFeeds() error {
	feeds, err := a.Models.FeedsModel.GetOrderedFeeds()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	for _, feed := range feeds {
		err := a.Models.FeedsModel.UpdatedLastFetched(feed.Url)
		if err != nil {
			log.Fatal(err)
		}

		rssFeed, err := fetchFeed(ctx, feed.Url)
		if err != nil{
			log.Fatal(err)
		}

		for i, _ := range rssFeed.Channel.Item {
			text, err := formatHtml(rssFeed.Channel.Item[i])
			if err != nil{
				return nil
			}

			fmt.Println(text)
		}	
	}

	return nil
}

func formatHtml(item RSSItem) (string, error) {
	
	var sb strings.Builder

	sb.WriteString(html.UnescapeString(item.Title))
	sb.WriteString(html.UnescapeString(item.Link))
	sb.WriteString(html.UnescapeString(item.Description))
	sb.WriteString(html.UnescapeString(item.PubDate))

	text, err := html2text.FromString(sb.String(), html2text.Options{TextOnly: true})
	if err != nil{
		return "", err
	}

	return text, nil
}