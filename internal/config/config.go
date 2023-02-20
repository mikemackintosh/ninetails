package config

// Default Colors
var DefaultColors = map[string]string{
	"PURPLE":    "38;5;129m",
	"PINK":      "38;5;162m",
	"RED":       "38;5;196m",
	"ORANGE":    "38;5;208m",
	"YELLOW":    "38;5;184m",
	"GREEN":     "38;5;154m",
	"BLUE":      "38;5;32m",
	"GREY":      "38;5;242m",
	"DARKGREY":  "38;5;239m",
	"LIGHTGREY": "38;5;249m",
	"BABYBLUE":  "38;5;123m",
	"LIGHTPINK": "38;5;212m",
	"WHITE":     "38;5;7m",
	"CLEAR":     "0m",
	"RESET":     "0m",
}

// Config is an instance of a configuration.
var Config *Configuration

// Configuration is the parsed structure of a configuration file.
type Configuration struct {
	Mappings Mappings          `yaml:"tails"`
	Colors   map[string]string `yaml:"colors,omitempty"`
}

// Replace is a helper method for looking through mappings for
// the given payload.
func Replace(s string) (string, bool) {
	var matched bool
	for _, m := range Config.Mappings {
		if m.re.MatchString(s) {
			s, matched = m.Replace(s)

			// If we get an exit on match, ditch it.
			// TODO: This should be refactored.
			if m.ExitOnMatch {
				return s, matched
			}
		}
	}

	return s, matched
}
