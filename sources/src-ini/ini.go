package src_ini

import (
	"congo"

	"github.com/go-ini/ini"
)

func New(path string) congo.Source {
	return &source{path}
}

type source struct {
	path string
}

func (s *source) Init(map[string]*congo.Setting) error {
	// Do nothing
	return nil
}

func (s *source) Load(settings map[string]*congo.Setting) error {
	cfg, err := ini.Load(s.path)
	if err != nil {
		return err
	}
	section, err := cfg.GetSection("")
	for key, setting := range settings {
		k, err := section.GetKey(key)
		if err != nil {
			return err
		}
		setting.Value.Set(k.Value())
	}
	return nil
}
