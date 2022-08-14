package config

import "github.com/sethvargo/go-githubactions"

// NewFromInputs creates a new Config from github actions configuration.
func NewFromInputs(action *githubactions.Action) *Config {
	return &Config{
		Name:     action.GetInput("name"),
		FilePath: action.GetInput("coverage-report"),
	}
}

// Config has the configuration to run go-coverage-report.
type Config struct {
	// run check name
	Name string
	// where coverage report is located.
	FilePath string
}
