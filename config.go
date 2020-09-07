// Read the JSON configuration
package hackerWecker

import (
	"encoding/json"
	"os"
)

type Config struct {
	Feeds          map[string]map[string][]string
	Music          []string
	NumberOfTracks int
	Shuffle        bool
}

func ReadConfig(configFile string) (Config, error) {
	// Read the config file encoded in JSON
	// Return a Config struct
	config := Config{}
	fh, err := os.Open(configFile)

	if err == nil {
		decoder := json.NewDecoder(fh)
		err = decoder.Decode(&config)
	}

	return config, err
}
