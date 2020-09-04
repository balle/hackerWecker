// Wake up a hacker by reading news and play some music or podcast
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/mmcdole/gofeed"
)

type Config struct {
	Feeds          map[string]map[string][]string
	Music          []string
	NumberOfTracks int
	Shuffle        bool
}

type fetchResult struct {
	Url     string
	Content string
}

func readConfig(feedFile string) Config {
	// Read the config file encoded in JSON
	// Return a Config struct
	fh, err := os.Open(feedFile)

	if err != nil {
		log.Fatal("Cannot read %s: %v", feedFile, err)
	}

	decoder := json.NewDecoder(fh)
	config := Config{}
	err = decoder.Decode(&config)

	if err != nil {
		log.Printf("Error decoding config: %v\n", err)
	}

	return config
}

func fetchUrl(url string, channel chan<- fetchResult) {
	// Get url and return fetchResult struct
	var result fetchResult
	result.Url = url

	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		speak(fmt.Sprintf("Error getting url %s: %v\n", url, err))
	} else {
		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			speak(fmt.Sprintf("Error reading url %s: %v\n", url, err))
			content = nil
		}

		result.Content = string(content)
	}

	channel <- result
}

func fetchFeeds(rssFeeds map[string]map[string][]string) map[string]string {
	// Fetch the contents of all feeds
	// Return a map of feed url as key and content as value
	content := make(map[string]string)
	channel := make(chan fetchResult)

	for url, _ := range rssFeeds {
		go fetchUrl(url, channel)
	}

	for i := 0; i < len(rssFeeds); i++ {
		result := <-channel

		if result.Content != "" {
			content[result.Url] = result.Content
		}
	}

	return content
}

func filterFeed(text string, metaData map[string][]string) bool {
	// Check if the feed should be read regarding to the include and exclude filters
	readFeed := true

	if metaData["exclude"] != nil {
		for i := 0; i < len(metaData["exclude"]); i++ {
			if strings.Contains(strings.ToLower(text), strings.ToLower(metaData["exclude"][i])) {
				readFeed = false
				break
			}
		}
	}

	if readFeed && metaData["include"] != nil {
		readFeed = false

		for i := 0; i < len(metaData["include"]); i++ {

			if strings.Contains(strings.ToLower(text), strings.ToLower(metaData["include"][i])) {
				readFeed = true
				break
			}
		}
	}

	return readFeed
}

func speak(text string) {
	// Filter the feed by using include / exclude metadata
	// Read the feed with help of a TTS tool
	ttsCmd := exec.Command("/usr/local/bin/espeak", "-a", "120", "-s", "150", "-v", "en-us")
	stdin, err := ttsCmd.StdinPipe()
	defer stdin.Close()

	if err != nil {
		log.Fatalf("Cannot pipe to espeak command: %v", err)
	}

	err = ttsCmd.Start()

	if err != nil {
		log.Fatalf("Cannot run espeak command: %v", err)
	}

	fmt.Fprintln(stdin, text)
	time.Sleep(time.Duration(len(text))*time.Millisecond*100 + 1*time.Second)
}

func readFeed(feed *gofeed.Feed, metaData map[string][]string) {
	// Check if a feed should be read regarding to the given metadata and if so read it
	fmt.Printf("////[ %s\n\n", feed.Title)
	speak(feed.Title)

	for _, item := range feed.Items {
		if filterFeed(item.Title, metaData) {
			fmt.Printf("Speak %s\n", item.Title)
			speak(item.Title)
		}

	}

	fmt.Println()
}

func remove(slice []string, i int) []string {
	// Remove element from slice of strings without preserving order
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func playMp3(filename string) {
	// Decode mp3 file and send it to the audio device
	fh, err := os.Open(filename)
	defer fh.Close()

	if err != nil {
		speak(fmt.Sprintf("Sorry cannot open %s: %v", filename, err))
	}

	decoder, err := mp3.NewDecoder(fh)

	if err != nil {
		speak(fmt.Sprintf("Sorry cannot decode %s: %v", filename, err))
		return
	}

	sound, err := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)

	if err != nil {
		speak(fmt.Sprintf("Sorry create sound interface: %v", err))
		return
	}

	defer sound.Close()

	player := sound.NewPlayer()
	defer player.Close()

	if _, err := io.Copy(player, decoder); err != nil {
		speak(fmt.Sprintf("Sorry cannot play %s: %v", filename, err))
		return
	}
}

func playMusic(musicDirs []string, numberOfTracks int, shuffle bool) {
	// Collect music files from given music dirs, if desired play randomly otherwise sequentially numberOfTracks
	var musicFiles []string

	for i := range musicDirs {
		fh, err := os.Open(musicDirs[i])
		defer fh.Close()

		dirList, err := fh.Readdir(-1)

		if err != nil {
			speak(fmt.Sprintf("Sorry cannot read directory %s: %v", musicDirs[i], err))
			continue
		}

		for x := range dirList {
			musicFiles = append(musicFiles, path.Join(musicDirs[i], dirList[x].Name()))
		}
	}

	if len(musicFiles) < numberOfTracks {
		numberOfTracks = len(musicFiles)
	}

	for i := 0; i < numberOfTracks; i++ {
		var playFile string

		if shuffle {
			rand.Seed(time.Now().Unix())
			x := rand.Intn(len(musicFiles))
			playFile = musicFiles[x]
			musicFiles = remove(musicFiles, x)
		} else {
			playFile = musicFiles[i]
		}

		fmt.Printf("Playing %s\n", playFile)
		playMp3(playFile)
	}
}

func main() {
	configFile := "hacker_wecker.json"

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	config := readConfig(configFile)
	contents := fetchFeeds(config.Feeds)
	parser := gofeed.NewParser()

	speak("Good morning, hacker!")
	playMusic(config.Music, config.NumberOfTracks, config.Shuffle)

	speak("Here are the news of the day.")

	for url, content := range contents {
		feed, err := parser.ParseString(content)

		if err != nil {
			speak(fmt.Sprintf("Cannot parse feed from %s: %v\n", url, err))
		} else {
			readFeed(feed, config.Feeds[url])
			time.Sleep(1 * time.Second)
		}
	}
}
