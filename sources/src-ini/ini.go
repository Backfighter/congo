package src_ini

import (
	"congo"

	"io"

	"github.com/go-ini/ini"
)

// New creates a new ini source which reads from
// given reader.
func New(reader io.ReadCloser) IniSource {
	return &iniSource{reader, ""}
}

// FromBytes creates a new ini source directly from
// the data that should be read.
func FromBytes(content []byte) IniSource {
	return &iniSource{content, ""}
}

// FromFile creates a new ini source which uses the
// file at given path to load the configuration.
func FromFile(path string) IniSource {
	return &iniSource{path, ""}
}

// IniSource a ini source uses input in ini-syntax
// to load settings.
type IniSource interface {
	congo.Source
	Section(name string) IniSource
}

type iniSource struct {
	source  interface{}
	section string
}

// Init initializes the ini source.
func (s *iniSource) Init(map[string]*congo.Setting) error {
	// Do nothing
	return nil
}

// Load loads the settings from input in ini-syntax.
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

// Section creates a sub-source that loads settings from a section
// of the ini input.
func (s *iniSource) Section(name string) IniSource {
	return &iniSource{s.source, name}
}
