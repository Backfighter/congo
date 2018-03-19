package src_ini

import (
	"congo"

	"io"

	"github.com/go-ini/ini"
)

func New(reader io.ReadCloser) IniSource {
	return &iniSource{reader, ""}
}

func FromBytes(content []byte) IniSource {
	return &iniSource{content, ""}
}

func FromFile(path string) IniSource {
	return &iniSource{path, ""}
}

type IniSource interface {
	congo.Source
	Section(name string) IniSource
}

type iniSource struct {
	source  interface{}
	section string
}

func (s *iniSource) Init(map[string]*congo.Setting) error {
	// Do nothing
	return nil
}

func (s *iniSource) Load(settings map[string]*congo.Setting) error {
	cfg, err := ini.Load(s.source)
	if err != nil {
		return err
	}
	section, err := cfg.GetSection(s.section)
	if err != nil {
		return err
	}
	for key, setting := range settings {
		if !section.HasKey(key) {
			continue
		}
		k, err := section.GetKey(key)
		if err != nil {
			return err
		}
		setting.Value.Set(k.Value())
	}
	return nil
}

func (s *iniSource) Section(name string) IniSource {
	return &iniSource{s.source, name}
}
