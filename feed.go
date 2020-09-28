// Handle RSS and Atom feeds
package hackerWecker

import (
	"fmt"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type Feed struct {
	Url   string
	Title string
	Items map[string]string
}

func FetchFeeds(feeds map[string]map[string][]string, outputChan chan<- Feed) {
	// Fetch the contents of all feeds, parse and filter them
	inputChan := make(chan fetchResult)
	maxAge := time.Now().AddDate(0, 0, config.MaxAgeOfFeedsInDays*-1)

	for url, _ := range feeds {
		go fetchUrl(url, inputChan)
	}

	for i := 0; i < len(feeds); i++ {
		input := <-inputChan

		go func(outputChan chan<- Feed) {
			feed, err := parseFeed(input.Url, input.Content)
			var result Feed

			if err == nil {
				result.Title = feed.Title
				result.Url = input.Url
				result.Items = make(map[string]string)

				for _, item := range feed.Items {
					if ((item.UpdatedParsed == nil && item.PublishedParsed != nil && item.PublishedParsed.Unix() > maxAge.Unix()) ||
						(item.UpdatedParsed != nil && item.UpdatedParsed.Unix() > maxAge.Unix())) &&
						filterFeed(input.Url, item.Title) {
						result.Items[item.Link] = item.Title
					} else if item.PublishedParsed == nil && item.UpdatedParsed == nil {
						LogInfo(fmt.Sprintf("skipping item without timestamp %s %s", feed.Title, item.Title))
					}
				}
			}

			outputChan <- result
		}(outputChan)
	}
}

func parseFeed(url string, content string) (*gofeed.Feed, error) {
	// Parse an RSS or Atom feed
	parser := gofeed.NewParser()

	feed, err := parser.ParseString(content)

	if err != nil {
		LogError(fmt.Sprintf("Cannot parse feed from %s: %v\n", url, err))
		feed = nil
	}

	return feed, err
}

func filterFeed(url, text string) bool {
	// Check if the feed should be read regarding to the include and exclude filters
	readFeed := true

	if config.Feeds[url]["exclude"] != nil {
		for i := 0; i < len(config.Feeds[url]["exclude"]); i++ {
			if strings.Contains(strings.ToLower(text), strings.ToLower(config.Feeds[url]["exclude"][i])) {
				readFeed = false
				break
			}
		}
	}

	if readFeed && config.Feeds[url]["include"] != nil {
		readFeed = false

		for i := 0; i < len(config.Feeds[url]["include"]); i++ {

			if strings.Contains(strings.ToLower(text), strings.ToLower(config.Feeds[url]["include"][i])) {
				readFeed = true
				break
			}
		}
	}

	return readFeed
}

func ReadFeed(feed Feed) int {
	// Read all feed items
	// Return number of read feed
	c := 0

	if len(feed.Items) > 0 {
		c++

		Speak(feed.Title)

		for _, title := range feed.Items {
			Speak(title)
		}

		fmt.Println()
	}

	return c
}

func GetFeeds() map[string]map[string][]string {
	return config.Feeds
}

func NumFeeds() int {
	return len(config.Feeds)
}

func GetPodcasts() map[string]map[string][]string {
	return config.Podcasts
}

func NumPodcasts() int {
	return len(config.Podcasts)
}
