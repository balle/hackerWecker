// Read the JSON configuration
package hackerWecker

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Feeds               map[string]map[string][]string
	FilterVars          map[string][]string
	MaxAgeOfFeedsInDays int
	Messages            map[string]string
	MusicDirs           []string
	NumberOfTracks      int
	OpenWeatherAPIKey   string
	Podcasts            map[string]map[string][]string
	Shuffle             bool
	TtsCmd              string
	TtsParams           string
	WeatherLocation     string
	WeatherUnit         string
	UserAgent           string
}

var config Config

func resolveVars(url string, metadata map[string][]string, option string) {
	// Resolve filter vars in include and exclude feed options
	if _, ok := metadata[option]; ok == true {
		for i := range metadata[option] {
			if strings.Contains(metadata[option][i], "var:") {
				varName := metadata[option][i]

				if _, ok := config.FilterVars[varName]; ok == true {
					config.Feeds[url][option] = remove(config.Feeds[url][option], i)

					for x := range config.FilterVars[varName] {
						config.Feeds[url][option] = append(config.Feeds[url][option], config.FilterVars[varName][x])
					}
				}
			}
		}
	}
}

func ReadConfig(configFile string) error {
	// Read the config file encoded in JSON
	// Resolve filter vars
	fh, err := os.Open(configFile)

	if err == nil {
		decoder := json.NewDecoder(fh)

		if err = decoder.Decode(&config); err != nil {
			err = fmt.Errorf("parsing JSON format in file %s: %v", configFile, err)
		}
	}

	for url, metadata := range config.Feeds {
		resolveVars(url, metadata, "include")
		resolveVars(url, metadata, "exclude")
	}

	return err
}

func GetMsg(msg string) string {
	return config.Messages[msg]
}
