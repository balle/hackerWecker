// URL handling stuff
package hackerWecker

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type FetchResult struct {
	Url     string
	Content string
	Error   error
}

func FetchUrl(url string, channel chan<- FetchResult) {
	// Get url and return fetchResult struct
	var result FetchResult
	result.Url = url
	rand.Seed(time.Now().Unix() * int64(os.Getpid()) / int64(os.Getppid()))

	for tries := 0; tries < 3; tries++ {
		fmt.Printf("Getting URL %s\n", url)
		resp, err := http.Get(url)
		defer resp.Body.Close()

		if err != nil {
			Speak(fmt.Sprintf("Error getting url %s: %v\n", url, err))
			result.Error = err
		} else {
			content, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				Speak(fmt.Sprintf("Error reading url %s: %v\n", url, err))
				result.Error = err
			} else {
				result.Content = string(content)
				break
			}
		}

		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	}

	channel <- result
}
