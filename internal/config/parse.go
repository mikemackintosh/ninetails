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

	// Set the instance config
	Config = &config

	// Check if default colors need to be set
	if len(config.Colors) == 0 {
		Config.Colors = DefaultColors
		return nil
	}

	// write back the default colors, skipping duplicates
	for color, code := range DefaultColors {
		if _, ok := config.Colors[color]; !ok {
			config.Colors[color] = code
		}
	}

	return nil
}
