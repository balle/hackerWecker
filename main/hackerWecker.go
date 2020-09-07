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
		log.Fatal("Cannot read %s: %v", configFile, err)
	}

	hackerWecker.Speak("Good morning, hacker!")

	channel := make(chan hackerWecker.FetchResult)
	go hackerWecker.FetchFeeds(config.Feeds, channel)
	hackerWecker.PlayMusic(config.Music, config.NumberOfTracks, config.Shuffle)

	hackerWecker.Speak("Here are the news of the day.")

	for i := 0; i < len(config.Feeds); i++ {
		result := <-channel

		if result.Error == nil {
			feed, err := hackerWecker.ParseFeed(result.Url, result.Content)

			if err == nil {
				hackerWecker.ReadFeed(feed, config.Feeds[result.Url])
				time.Sleep(1 * time.Second)
			}
		}
	}
}
