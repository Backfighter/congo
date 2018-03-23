package env

import (
	"os"

	"strings"

	"regexp"

	"gitlab.com/silentteacup/congo"
)

/*
Copyright (c) 2018 Peter Werner. All rights reserved.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
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
