// Wake up a hacker by reading news and play some music or podcast
package main

import (
	"log"
	"os"
	"time"

	"github.com/balle/hackerWecker"
)

func main() {
	configFile := "hackerWecker.json"

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	config, err := hackerWecker.ReadConfig(configFile)

	if err != nil {
		log.Fatalf("Cannot read %s: %v", configFile, err)
	}

	hackerWecker.Speak("Good morning, hacker!")

	channel := make(chan hackerWecker.FeedResult)
	go hackerWecker.FetchFeeds(config, channel)
	hackerWecker.PlayMusic(config)

	hackerWecker.Speak("Here are the news of the day.")

	for i := 0; i < len(config.Feeds); i++ {
		feed := <-channel

		hackerWecker.ReadFeed(feed, config.Feeds[feed.Url])
		time.Sleep(1 * time.Second)
	}
}
