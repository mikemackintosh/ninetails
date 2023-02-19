package config

import (
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// Parse will parse the configuration.
func Parse(f string) error {
	file, err := os.ReadFile(f)
	if err != nil {
		return err
	}

	var config Configuration
	if err = yaml.Unmarshal(file, &config); err != nil {
		return err
	}

	for _, m := range config.Mappings {
		r, err := regexp.Compile(m.Search)
		if err != nil {
			log.Printf("Skipping invalid configuration for: %s", m.Search)
			continue
		}
		m.re = r
	}

	Config = &config
	return nil
}
