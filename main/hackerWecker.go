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

	chanFeeds := make(chan hackerWecker.Feed)
	go hackerWecker.FetchFeeds(config.Feeds, chanFeeds)

	chanPodcasts := make(chan hackerWecker.Feed)
	go hackerWecker.FetchFeeds(config.Podcasts, chanPodcasts)

	hackerWecker.PlayMusic()
	hackerWecker.Speak("Here are the news of the day.")

	for i := 0; i < len(config.Feeds); i++ {
		feed := <-chanFeeds

		hackerWecker.ReadFeed(feed)
		time.Sleep(1 * time.Second)
	}

	for i := 0; i < len(config.Podcasts); i++ {
		feed := <-chanPodcasts

		hackerWecker.PlayPodcast(feed)
		time.Sleep(1 * time.Second)
	}
}
