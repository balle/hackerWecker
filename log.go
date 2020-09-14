package hackerWecker

import (
	"fmt"
	"log"
)

func LogInfo(msg string) {
	fmt.Println(msg)
}

func LogError(msg string) {
	log.Println(msg)
	Speak(msg)
}

func LogFatal(msg string) {
	log.Fatal(msg)
}
