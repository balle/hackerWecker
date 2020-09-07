// URL handling stuff
package hackerWecker

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type fetchResult struct {
	Url     string
	Content string
}

func FetchUrl(url string, channel chan<- fetchResult) {
	// Get url and return fetchResult struct
	var result fetchResult
	result.Url = url

	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		Speak(fmt.Sprintf("Error getting url %s: %v\n", url, err))
	} else {
		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			Speak(fmt.Sprintf("Error reading url %s: %v\n", url, err))
			content = nil
		}

		result.Content = string(content)
	}

	channel <- result
}
