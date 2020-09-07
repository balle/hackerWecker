// Handle RSS feeds
package hackerWecker

import (
	"fmt"
	"strings"

	"github.com/mmcdole/gofeed"
)

func FetchFeeds(rssFeeds map[string]map[string][]string) map[string]string {
	// Fetch the contents of all feeds
	// Return a map of feed url as key and content as value
	content := make(map[string]string)
	channel := make(chan fetchResult)

	for url, _ := range rssFeeds {
		go FetchUrl(url, channel)
	}

	for i := 0; i < len(rssFeeds); i++ {
		result := <-channel

		if result.Content != "" {
			content[result.Url] = result.Content
		}
	}

	return content
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

func ReadFeed(feed *gofeed.Feed, metaData map[string][]string) {
	// Check if a feed should be read regarding to the given metadata and if so read it
	fmt.Printf("////[ %s\n\n", feed.Title)
	Speak(feed.Title)

	for _, item := range feed.Items {
		if FilterFeed(item.Title, metaData) {
			fmt.Printf("Speak %s\n", item.Title)
			Speak(item.Title)
		}

	}

	fmt.Println()
}
