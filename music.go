// Play some music (mp3)
package hackerWecker

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

func remove(slice []string, i int) []string {
	// Remove element from slice of strings without preserving order
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func playMp3FromFilehandle(fh io.Reader, filename string) {
	// Decode MP3 from filehandle and send it to the audio device
	decoder, err := mp3.NewDecoder(fh)

	if err != nil {
		LogError(fmt.Sprintf("Sorry cannot decode %s: %v", filename, err))
		return
	}

	sound, err := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)

	if err != nil {
		LogError(fmt.Sprintf("Sorry create sound interface: %v", err))
		return
	}

	defer sound.Close()

	player := sound.NewPlayer()
	defer player.Close()

	if _, err := io.Copy(player, decoder); err != nil {
		LogError(fmt.Sprintf("Sorry cannot play %s: %v", filename, err))
		return
	}
}

func playMp3(filename string) {
	// Open mp3 file, decode it and send it to the audio device
	fh, err := os.Open(filename)
	defer fh.Close()

	if err != nil {
		LogError(fmt.Sprintf("Sorry cannot open %s: %v", filename, err))
		return
	}

	playMp3FromFilehandle(fh, filename)
}

func PlayMusic() {
	// Collect music files from given music dirs, if desired play randomly otherwise sequentially numberOfTracks
	var musicFiles []string
	var numberOfTracks int

	for i := range config.MusicDirs {
		fh, err := os.Open(config.MusicDirs[i])
		defer fh.Close()

		dirList, err := fh.Readdir(-1)

		if err != nil {
			LogError(fmt.Sprintf("Sorry cannot read directory %s: %v", config.MusicDirs[i], err))
			continue
		}

		for x := range dirList {
			musicFiles = append(musicFiles, path.Join(config.MusicDirs[i], dirList[x].Name()))
		}
	}

	if len(musicFiles) < config.NumberOfTracks {
		numberOfTracks = len(musicFiles)
	} else {
		numberOfTracks = config.NumberOfTracks
	}

	for i := 0; i < numberOfTracks; i++ {
		var playFile string

		if config.Shuffle {
			rand.Seed(time.Now().Unix() * int64(os.Getpid()) / int64(os.Getppid()))
			x := rand.Intn(len(musicFiles))
			playFile = musicFiles[x]
			musicFiles = remove(musicFiles, x)
		} else {
			playFile = musicFiles[i]
		}

		LogInfo("Playing " + playFile)
		playMp3(playFile)
	}
}

func PlayPodcast(feed Feed) {
	// Stream Podcast mp3 to audio device
	for url, title := range feed.Items {
		LogInfo("Playing podcast " + title)
		Speak(feed.Title)

		client, req := initWebReq(url)
		resp, err := client.Do(req)
		defer resp.Body.Close()

		if err == nil {
			playMp3FromFilehandle(resp.Body, "podcast "+title)
		}
	}

}
