package src_env

import (
	"congo"
	"os"
)

// New creates a new environment source. Which directly
// loads settings from environment variables.
func New() congo.Source {
	return &source{}
}

type source struct{}

// Inits initializes this source
func (s *source) Init(map[string]*congo.Setting) error {
	// Do nothing
	return nil
}

// Load loads settings from environment variables.
func (s *source) Load(settings map[string]*congo.Setting) error {
	for key, setting := range settings {
		value, ok := os.LookupEnv(key)
		if ok {
			setting.Value.Set(value)
		}
	}
	return nil
}
