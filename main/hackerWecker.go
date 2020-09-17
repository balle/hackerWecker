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

	err := hackerWecker.ReadConfig(configFile)

	if err != nil {
		hackerWecker.LogFatal(fmt.Sprintf("Cannot read %s: %v", configFile, err))
	}

	hackerWecker.Speak(hackerWecker.GetMsg("welcome"))

	chanFeeds := make(chan hackerWecker.Feed)
	go hackerWecker.FetchFeeds(hackerWecker.GetFeeds(), chanFeeds)

	chanPodcasts := make(chan hackerWecker.Feed)
	go hackerWecker.FetchFeeds(hackerWecker.GetPodcasts(), chanPodcasts)

	//hackerWecker.PlayMusic()
	hackerWecker.Speak(hackerWecker.GetMsg("news"))

	for i := 0; i < len(hackerWecker.GetFeeds()); i++ {
		feed := <-chanFeeds

		hackerWecker.ReadFeed(feed)
		time.Sleep(1 * time.Second)
	}

	for i := 0; i < len(hackerWecker.GetPodcasts()); i++ {
		feed := <-chanPodcasts

		hackerWecker.PlayPodcast(feed)
		time.Sleep(1 * time.Second)
	}

	hackerWecker.Speak(hackerWecker.GetMsg("finished"))
}
