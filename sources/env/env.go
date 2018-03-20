package env

import (
	"os"

	"strings"

	"regexp"

	"gitlab.com/silentteacup/congo"
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

// New creates a new environment source. Which directly
// loads settings from environment variables.
func New() Source {
	return &source{IdenticalTranslator}
}

// Translator is used to translate the settings names to more conventional
// ones.
type Translator func(string) []string

// IdenticalTranslator is a translator that always returns the
// key in a string array of size 1.
func IdenticalTranslator(key string) []string {
	return []string{key}
}

var sdtTranslatorFilter = regexp.MustCompile("[^A-Z0-9_]")

// PrefixSdtTranslator returns a translator that translates the key string
// into a environment variable form. The given prefix will be added in front
// of the key and after that spaces and hyphens will be converted to '_' and the string
// will be converted to uppercase. Finally all characters except A-Z and 0-9 or _
// will be removed.
func PrefixSdtTranslator(prefix string) Translator {
	return func(s string) []string {
		s = strings.ToUpper(strings.Replace(prefix+s, " ", "_", -1))
		s = strings.Replace(s, "-", "_", -1)
		return []string{sdtTranslatorFilter.ReplaceAllString(s, "")}
	}
}

// Source is a source that gathers the settings from environment variables.
type Source interface {
	congo.Source
	// WithTranslator add a translator function that translates a
	// defined name for a setting into its alternative representations.
	// When looking for the environment variable all of these representations
	// will be considered for the key.
	//
	// The alternative representations are ordered representations given first will
	// be preferred over others.
	WithTranslator(t Translator) Source
}

type source struct {
	translator Translator
}

// WithTranslator add a translator function that translates a
// defined name for a setting into its alternative representations.
// When looking for the environment variable all of these representations
// will be considered for the key.
//
// The alternative representations are ordered representations given first will
// be preferred over others.
func (s *source) WithTranslator(t Translator) Source {
	s.translator = t
	return s
}

// Inits initializes this source
func (s *source) Init(map[string]*congo.Setting) error {
	// Do nothing
	return nil
}

// Load loads settings from environment variables.
func (s *source) Load(settings map[string]*congo.Setting) error {
	for key, setting := range settings {
		for _, alternative := range s.translator(key) {
			value, ok := os.LookupEnv(alternative)
			if ok {
				setting.Value.Set(value)
				break
			}
		}
	}
	return nil
}
