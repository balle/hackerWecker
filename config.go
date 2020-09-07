// Read the JSON configuration
package hackerWecker

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Feeds          map[string]map[string][]string
	Music          []string
	NumberOfTracks int
	Shuffle        bool
}

func ReadConfig(feedFile string) Config {
	// Read the config file encoded in JSON
	// Return a Config struct
	fh, err := os.Open(feedFile)

	if err != nil {
		log.Fatal("Cannot read %s: %v", feedFile, err)
	}

	decoder := json.NewDecoder(fh)
	config := Config{}
	err = decoder.Decode(&config)

	if err != nil {
		log.Printf("Error decoding config: %v\n", err)
	}

	return config
}
