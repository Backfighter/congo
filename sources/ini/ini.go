package ini

import (
	"gitlab.com/silentteacup/congo"

	"io"

	"fmt"

	"github.com/go-ini/ini"
)

/*
Copyright (c) 2018 Peter Werner. All rights reserved.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

// New creates a new ini source which reads from
// given reader.
func New(reader io.ReadCloser) Source {
	return createSource(reader)
}

// FromBytes creates a new ini source directly from
// the data that should be read.
func FromBytes(content []byte) Source {
	return createSource(content)
}

// FromFile creates a new ini source which uses the
// file at given path to load the configuration.
func FromFile(path string) Source {
	return createSource(path)
}

// createSource creates the ini source with default values
// using given source as source for the ini-file.
func createSource(source interface{}) Source {
	return &iniSource{source, "", true, nil}
}

// Source a ini source uses input in ini-syntax
// to load settings.
type Source interface {
	congo.Source
	Section(name string) Source
	WriteDefaults(w io.Writer) error
	SetLooseLoad(loose bool) Source
}

type iniSource struct {
	source    interface{}
	section   string
	looseLoad bool
	defaults  map[string]*congo.Setting
}

// Init initializes the ini source.
func (s *iniSource) Init(settings map[string]*congo.Setting) error {
	s.defaults = settings
	return nil
}

// loadIni loads the ini file in the appropriate way.
func (s *iniSource) loadIni() (cfg *ini.File, err error) {
	if s.looseLoad {
		cfg, err = ini.LooseLoad(s.source)
	} else {
		cfg, err = ini.Load(s.source)
	}
	return cfg, err
}

// Load loads the settings from input in ini-syntax.
func (s *iniSource) Load(settings map[string]*congo.Setting) error {
	cfg, err := s.loadIni()
	if err != nil {
		return fmt.Errorf("ini-source: couldn't load the ini-file because: %s", err)
	}
	section, err := cfg.GetSection(s.section)
	if err != nil {
		// Section doesn't exist
		// We simply don't load the section and use the defaults
		return nil
	}
	for key, setting := range settings {
		if !section.HasKey(key) {
			continue
		}
		k, err := section.GetKey(key)
		if err != nil {
			// Key exists but can't get key
			// Return error
			return err
		}
		if err := setting.Value.Set(k.Value()); err != nil {
			return fmt.Errorf("ini-source: couldn't read setting %q "+
				"in section %q: %s", key, s.section, err)
		}
	}
	return nil
}

// WriteDefaults writes the default settings to given writer.
// If an error occurs nothing will be written.
func (s *iniSource) WriteDefaults(w io.Writer) (err error) {
	cfg := ini.Empty()
	section := cfg.Section(s.section)
	for name, setting := range s.defaults {
		k := section.Key(name)
		k.Comment = setting.Usage
		k.SetValue(setting.DefValue)
	}
	_, err = cfg.WriteTo(w)
	return err
}

// SetLooseLoad sets whether this source should complain if the file
// doesn't exist. Default is true.
func (s *iniSource) SetLooseLoad(loose bool) Source {
	s.looseLoad = loose
	return s
}

// Section creates a sub-source that loads settings from a section
// of the ini input.
func (s *iniSource) Section(name string) Source {
	return &iniSource{
		s.source,
		name,
		s.looseLoad,
		s.defaults,
	}
}
