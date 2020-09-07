// Wake up a hacker by reading news and play some music or podcast
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/balle/hackerWecker"
	"github.com/mmcdole/gofeed"
)

func main() {
	configFile := "hacker_wecker.json"

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	config := hackerWecker.ReadConfig(configFile)
	contents := hackerWecker.FetchFeeds(config.Feeds)
	parser := gofeed.NewParser()

	hackerWecker.Speak("Good morning, hacker!")
	hackerWecker.PlayMusic(config.Music, config.NumberOfTracks, config.Shuffle)

	hackerWecker.Speak("Here are the news of the day.")

	for url, content := range contents {
		feed, err := parser.ParseString(content)

		if err != nil {
			hackerWecker.Speak(fmt.Sprintf("Cannot parse feed from %s: %v\n", url, err))
		} else {
			hackerWecker.ReadFeed(feed, config.Feeds[url])
			time.Sleep(1 * time.Second)
		}
	}
}
