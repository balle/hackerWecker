// Wake up a hacker by reading news and play some music or podcast
package main

import (
	"fmt"
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
		hackerWecker.LogFatal(fmt.Sprintf("Cannot read %s: %v", configFile, err))
	}

	hackerWecker.Speak(hackerWecker.GetMsg("welcome"))

	chanFeeds := make(chan hackerWecker.Feed)
	go hackerWecker.FetchFeeds(config.Feeds, chanFeeds)

	chanPodcasts := make(chan hackerWecker.Feed)
	go hackerWecker.FetchFeeds(config.Podcasts, chanPodcasts)

	hackerWecker.PlayMusic()
	hackerWecker.Speak(hackerWecker.GetMsg("news"))

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

	hackerWecker.Speak(hackerWecker.GetMsg("finished"))
}
