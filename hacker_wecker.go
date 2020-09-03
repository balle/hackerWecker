// Wake up a hacker by reading news and play some music or podcast
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type Config struct {
	Feeds map[string]map[string][]string
}

type fetchResult struct {
	Url     string
	Content string
}

func readConfig(feedFile string) Config {
	// Read the config file encoded in JSON
	// Return a Config struct
	fh, err := os.Open(feedFile)

	if err != nil {
		log.Fatal("Cannot read %s: %v", feedFile, err)
	}

	decoder := json.NewDecoder(fh)
	config := Config{}
	err = decoder.Decode(&config)

	if err != nil {
		log.Printf("Error decoding config: %v\n", err)
	}

	return config
}

func fetchUrl(url string, channel chan<- fetchResult) {
	// Get url and return fetchResult struct
	var result fetchResult
	result.Url = url

	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting url %s: %v\n", url, err)
	} else {
		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading url %s: %v\n", url, err)
			content = nil
		}

		result.Content = string(content)
	}

	channel <- result
}

func fetchFeeds(rssFeeds map[string]map[string][]string) map[string]string {
	// Fetch the contents of all feeds
	// Return a map of feed url as key and content as value
	content := make(map[string]string)
	channel := make(chan fetchResult)

	for url, _ := range rssFeeds {
		go fetchUrl(url, channel)
	}

	for i := 0; i < len(rssFeeds); i++ {
		result := <-channel

		if result.Content != "" {
			content[result.Url] = result.Content
		}
	}

	return content
}

func filterFeed(text string, metaData map[string][]string) bool {
	readFeed := true

	// Check if the feed should be read regarding to the include and exclude filters
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

func speak(text string) {
	// Filter the feed by using include / exclude metadata
	// Read the feed with help of a TTS tool
	ttsCmd := exec.Command("/usr/local/bin/espeak", "-a", "120", "-s", "150", "-v", "en-us")
	stdin, err := ttsCmd.StdinPipe()
	defer stdin.Close()

	if err != nil {
		log.Fatalf("Cannot pipe to espeak command: %v", err)
	}

	err = ttsCmd.Start()

	if err != nil {
		log.Fatalf("Cannot run espeak command: %v", err)
	}

	fmt.Fprintln(stdin, text)
	time.Sleep(time.Duration(len(text))*time.Millisecond*100 + 1*time.Second)
}

func readFeed(feed *gofeed.Feed, metaData map[string][]string) {
	fmt.Printf("////[ %s\n\n", feed.Title)
	speak(feed.Title)

	for _, item := range feed.Items {
		if filterFeed(item.Title, metaData) {
			fmt.Printf("Speak %s\n", item.Title)
			speak(item.Title)
		}

	}

	fmt.Println()
}

func main() {
	configFile := "hacker_wecker.json"

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	config := readConfig(configFile)
	contents := fetchFeeds(config.Feeds)
	parser := gofeed.NewParser()

	speak("Good morning. Here are the news of the day.")

	for url, content := range contents {
		feed, err := parser.ParseString(content)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot parse feed from %s: %v\n", url, err)
		} else {
			readFeed(feed, config.Feeds[url])
			time.Sleep(1 * time.Second)
		}
	}
}
