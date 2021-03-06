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

type fetchResult struct {
	Url     string
	Content string
	Error   error
}

func initWebReq(url string) (*http.Client, *http.Request) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", config.UserAgent)

	return client, req
}

func fetchUrl(url string, channel chan<- fetchResult) {
	// Get url and return fetchResult struct
	var result fetchResult
	result.Url = url

	rand.Seed(time.Now().Unix() * int64(os.Getpid()) / int64(os.Getppid()))

	for tries := 0; tries < 3; tries++ {
		fmt.Printf("Getting URL %s\n", url)
		client, req := initWebReq(url)
		resp, err := client.Do(req)

		if err != nil {
			Speak(fmt.Sprintf("Error getting url %s: %v\n", url, err))
			result.Error = err
		} else {
			defer resp.Body.Close()
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
