package src_env

import (
	"congo"
	"os"
)

func New() congo.Source {
	return &source{}
}

type source struct{}

func (s *source) Init(map[string]*congo.Setting) error {
	// Do nothing
	return nil
}

func (s *source) Load(settings map[string]*congo.Setting) error {
	for key, setting := range settings {
		value, ok := os.LookupEnv(key)
		if ok {
			setting.Value.Set(value)
		}
	}
	return nil
}
