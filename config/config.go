package config

// Config is an instance of a configuration
var Config = &Configuration{}

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
			return m.Replace(s)
		}
	}

	return s, matched
}
