// Text to speech
package hackerWecker

import (
	"fmt"
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
		LogFatal(fmt.Sprintf("Cannot pipe to espeak command: %v", err))
	}

	err = ttsCmd.Start()

	if err != nil {
		LogFatal(fmt.Sprintf("Cannot run espeak command: %v", err))
	}

	fmt.Fprintln(stdin, text)
	LogInfo(text)
	time.Sleep(time.Duration(len(text))*time.Millisecond*100 + 1*time.Second)
}
