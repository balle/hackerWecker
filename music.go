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

func PlayMp3(filename string) {
	// Decode mp3 file and send it to the audio device
	fh, err := os.Open(filename)
	defer fh.Close()

	if err != nil {
		Speak(fmt.Sprintf("Sorry cannot open %s: %v", filename, err))
	}

	decoder, err := mp3.NewDecoder(fh)

	if err != nil {
		Speak(fmt.Sprintf("Sorry cannot decode %s: %v", filename, err))
		return
	}

	sound, err := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)

	if err != nil {
		Speak(fmt.Sprintf("Sorry create sound interface: %v", err))
		return
	}

	defer sound.Close()

	player := sound.NewPlayer()
	defer player.Close()

	if _, err := io.Copy(player, decoder); err != nil {
		Speak(fmt.Sprintf("Sorry cannot play %s: %v", filename, err))
		return
	}
}

func PlayMusic(musicDirs []string, numberOfTracks int, shuffle bool) {
	// Collect music files from given music dirs, if desired play randomly otherwise sequentially numberOfTracks
	var musicFiles []string

	for i := range musicDirs {
		fh, err := os.Open(musicDirs[i])
		defer fh.Close()

		dirList, err := fh.Readdir(-1)

		if err != nil {
			Speak(fmt.Sprintf("Sorry cannot read directory %s: %v", musicDirs[i], err))
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
		PlayMp3(playFile)
	}
}
