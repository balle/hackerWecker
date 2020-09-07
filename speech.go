// Text to speech
package hackerWecker

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

func Speak(text string) {
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
