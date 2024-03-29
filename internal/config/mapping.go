package config

import (
	"fmt"
	"regexp"
	"strings"
)

// Mappings is a slice of Mapping pointers.
type Mappings []*Mapping

// Mapping is a struct of match instructions.
//   - Search is the string you want to match on
//   - Color is the color to highlight the matching line
//   - Format is the regex replacement format you want to output
//   - ExitOnMatch determines if we should stop processing after
//     the rule is matched
type Mapping struct {
	Search      string  `yaml:"search"`
	Color       *string `yaml:"color"`
	Format      *string `yaml:"format"`
	ExitOnMatch bool    `yaml:"exit_on_match"`

	re *regexp.Regexp
}

// Replace will replace the string with the formatted string
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
