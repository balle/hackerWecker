// Handle RSS feeds
package hackerWecker

import (
	"fmt"
	"strings"

	"github.com/mmcdole/gofeed"
)

type FeedResult struct {
	Url   string
	Title string
	Items []string
}

func FetchFeeds(config Config, outputChan chan<- FeedResult) {
	// Fetch the contents of all feeds
	// Return a map of feed url as key and content as value
	inputChan := make(chan FetchResult)

	for url, _ := range config.Feeds {
		go FetchUrl(url, inputChan)
	}

	for i := 0; i < len(config.Feeds); i++ {
		input := <-inputChan
		feed, err := ParseFeed(input.Url, input.Content)
		var result FeedResult

		if err == nil {
			result.Url = input.Url

			for _, item := range feed.Items {
				if FilterFeed(item.Title, config.Feeds[input.Url]) {
					result.Items = append(result.Items, item.Title)
				}
			}
		}

		outputChan <- result
	}
}

func ParseFeed(url string, content string) (*gofeed.Feed, error) {
	// Parse an RSS feed
	parser := gofeed.NewParser()

	feed, err := parser.ParseString(content)

	if err != nil {
		Speak(fmt.Sprintf("Cannot parse feed from %s: %v\n", url, err))
		feed = nil
	}

	return feed, err
}

func FilterFeed(text string, metaData map[string][]string) bool {
	// Check if the feed should be read regarding to the include and exclude filters
	readFeed := true

	if metaData["exclude"] != nil {
		for i := 0; i < len(metaData["exclude"]); i++ {
			if strings.Contains(strings.ToLower(text), strings.ToLower(metaData["exclude"][i])) {
				readFeed = false
				break
			}
		}
	}

	if readFeed && metaData["include"] != nil {
		readFeed = false

		for i := 0; i < len(metaData["include"]); i++ {

			if strings.Contains(strings.ToLower(text), strings.ToLower(metaData["include"][i])) {
				readFeed = true
				break
			}
		}
	}

	return readFeed
}

func ReadFeed(feed FeedResult, metaData map[string][]string) {
	// Check if a feed should be read regarding to the given metadata and if so read it
	fmt.Printf("////[ %s\n\n", feed.Title)
	Speak(feed.Title)

	for _, item := range feed.Items {
		if FilterFeed(item, metaData) {
			fmt.Printf("Speak %s\n", item)
			Speak(item)
		}

	}

	fmt.Println()
}
