// Wake up a hacker by reading news and play some music or podcast
package main

import (
	"os"
	"time"

	"github.com/balle/hackerWecker"
)

func main() {
	configFile := "hackerWecker.json"

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	config := hackerWecker.ReadConfig(configFile)
	contents := hackerWecker.FetchFeeds(config.Feeds)

	hackerWecker.Speak("Good morning, hacker!")
	hackerWecker.PlayMusic(config.Music, config.NumberOfTracks, config.Shuffle)

	hackerWecker.Speak("Here are the news of the day.")

	for url, content := range contents {
		feed := hackerWecker.ParseFeed(url, content)

		if feed != nil {
			hackerWecker.ReadFeed(feed, config.Feeds[url])
			time.Sleep(1 * time.Second)
		}
	}
}
