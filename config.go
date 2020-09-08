// Read the JSON configuration
package hackerWecker

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Feeds               map[string]map[string][]string
	MaxAgeOfFeedsInDays int
	MusicDirs           []string
	NumberOfTracks      int
	Shuffle             bool
	TtsCmd              string
	TtsParams           string
}

var config Config

func ReadConfig(configFile string) (Config, error) {
	// Read the config file encoded in JSON
	// Return a Config struct
	fh, err := os.Open(configFile)

	if err == nil {
		decoder := json.NewDecoder(fh)

		if err = decoder.Decode(&config); err != nil {
			err = fmt.Errorf("parsing JSON format in file %s: %v", configFile, err)
		}
	}

	return config, err
}
