// Wake up a hacker by reading news and play some music or podcast
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/balle/hackerWecker"
)

func main() {
	var configFile = flag.String("config", "hackerWecker.json", "config file")
	readFeeds := 0

	flag.Parse()
	err := hackerWecker.ReadConfig(*configFile)

	if err != nil {
		hackerWecker.LogFatal(fmt.Sprintf("Cannot read %s: %v", *configFile, err))
	}

	hackerWecker.Speak(hackerWecker.GetMsg("welcome"))

	chanFeeds := make(chan hackerWecker.Feed)
	go hackerWecker.FetchFeeds(hackerWecker.GetFeeds(), chanFeeds)

	chanPodcasts := make(chan hackerWecker.Feed)
	go hackerWecker.FetchFeeds(hackerWecker.GetPodcasts(), chanPodcasts)

	hackerWecker.PlayMusic()
	hackerWecker.Speak(hackerWecker.GetMsg("news"))

	for i := 0; i < len(hackerWecker.GetFeeds()); i++ {
		readFeeds += hackerWecker.ReadFeed(<-chanFeeds)
		time.Sleep(1 * time.Second)
	}

	if readFeeds == 0 {
		hackerWecker.Speak(hackerWecker.GetMsg("nonews"))
	}

	hackerWecker.Speak(hackerWecker.GetMsg("podcasts"))

	for i := 0; i < len(hackerWecker.GetPodcasts()); i++ {
		hackerWecker.PlayPodcast(<-chanPodcasts)
	}

	hackerWecker.Speak(hackerWecker.GetMsg("finished"))
}
