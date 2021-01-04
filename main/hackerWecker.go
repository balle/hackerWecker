// Wake up a hacker by reading news and play some music and podcast
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/balle/hackerWecker"
)

func main() {
	var configFile = flag.String("config", "hackerWecker.json", "config file")
	var feeds []hackerWecker.Feed

	flag.Parse()
	err := hackerWecker.ReadConfig(*configFile)

	if err != nil {
		hackerWecker.LogFatal(fmt.Sprintf("Cannot read %s: %v", *configFile, err))
	}

	hackerWecker.SetupMixer()

	chanFeeds := make(chan hackerWecker.Feed, hackerWecker.NumFeeds())
	chanPodcasts := make(chan hackerWecker.Feed, hackerWecker.NumPodcasts())

	hackerWecker.Speak(hackerWecker.GetMsg("welcome"))

	go hackerWecker.FetchFeeds(hackerWecker.GetFeeds(), chanFeeds)
	go hackerWecker.FetchFeeds(hackerWecker.GetPodcasts(), chanPodcasts)

	hackerWecker.PlayMusic()

	for i := 0; i < hackerWecker.NumFeeds(); i++ {
		feeds = append(feeds, <-chanFeeds)
	}

	if len(feeds) == 0 {
		hackerWecker.Speak(hackerWecker.GetMsg("nonews"))
	} else {
		hackerWecker.Speak(hackerWecker.GetMsg("news"))

		for _, feed := range feeds {
			hackerWecker.ReadFeed(feed)
			time.Sleep(1 * time.Second)
		}
	}

	hackerWecker.Speak(hackerWecker.GetMsg("weather"))
	hackerWecker.ReadWeather()

	hackerWecker.Speak(hackerWecker.GetMsg("podcasts"))

	for i := 0; i < hackerWecker.NumPodcasts(); i++ {
		hackerWecker.PlayPodcast(<-chanPodcasts)
	}

	hackerWecker.Speak(hackerWecker.GetMsg("finished"))
}
