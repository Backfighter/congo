package src_ini

import (
	"gitlab.com/silentteacup/congo"

	"io"

	"github.com/go-ini/ini"
)

/*
Copyright (c) 2018 Peter Werner. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

* Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
* Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

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
