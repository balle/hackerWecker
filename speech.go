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
	ttsCmd := exec.Command(config.TtsCmd, config.TtsParams)
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
	fmt.Println(text)
	time.Sleep(time.Duration(len(text))*time.Millisecond*100 + 1*time.Second)
}
