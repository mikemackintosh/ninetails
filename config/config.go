package config

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var Config = &Configuration{}

type Mappings []*Mapping
type Mapping struct {
	Search string  `yaml:"search"`
	Color  *string `yaml:"color"`
	Format *string `yaml:"format"`

	re *regexp.Regexp
}

func (m *Mapping) Replace(s string) (string, bool) {
	if m.Color != nil {
		if c, ok := Config.Colors[*m.Color]; ok {
			return fmt.Sprintf("\033[%s%s\033[0m", c, s), true
		}
	}

	if m.Format != nil {
		var format = *m.Format
		for k, v := range Config.Colors {
			format = strings.Replace(format, "\\"+k, "\033["+v, -1)
		}
		v := m.re.ReplaceAllString(s, format+"\033[0m")
		return v, true
	}

	return "", false
}

type Color string
type Configuration struct {
	Mappings Mappings          `yaml:"tails"`
	Colors   map[string]string `yaml:"colors,omitempty"`
}

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

func Replace(s string) (string, bool) {
	var matched bool
	for _, m := range Config.Mappings {
		if m.re.MatchString(s) {
			return m.Replace(s)
		}
	}

	return s, matched
}
