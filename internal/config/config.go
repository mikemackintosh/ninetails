package config

// Config is an instance of a configuration.
var Config = &Configuration{}

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

			// If we are not combining the results,
			// we should return early
			if m.ExitOnMatch {
				return s, matched
			}
		}
	}

	return s, matched
}
