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

	contents := hackerWecker.FetchFeeds(config.Feeds)

	hackerWecker.Speak("Good morning, hacker!")
	hackerWecker.PlayMusic(config.Music, config.NumberOfTracks, config.Shuffle)

	hackerWecker.Speak("Here are the news of the day.")

	for url, content := range contents {
		feed, err := hackerWecker.ParseFeed(url, content)

		if err == nil {
			hackerWecker.ReadFeed(feed, config.Feeds[url])
			time.Sleep(1 * time.Second)
		}
	}
}
