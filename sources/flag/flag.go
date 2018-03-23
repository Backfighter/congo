package flag

import (
	"flag"
	"os"

	"gitlab.com/silentteacup/congo"
)

/*
Copyright (c) 2018 Peter Werner. All rights reserved.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

// New creates a new flag source using the standard command
// line FlagSet(flag.CommandLine) and arguments from the command line.
func New() congo.Source {
	return FromFlagSet(flag.CommandLine, standardLoader)
}

// ArgLoader is used to load arguments when parsing flags.
type ArgLoader func() []string

// FromFlagSet creates a new flag source using a custom FlagSet and
// argument loader. The argument loader specifies how arguments are loaded
// when the flags are parsed.
func FromFlagSet(set *flag.FlagSet, loader ArgLoader) congo.Source {
	return &source{set, loader}
}

// standardLoader loads the commandline arguments
func standardLoader() []string {
	return os.Args[1:]
}

type source struct {
	set *flag.FlagSet
	ArgLoader
}

// Init registers the flags for this source
func (s *source) Init(settings map[string]*congo.Setting) error {
	for key, setting := range settings {
		s.set.Var(setting.Value, key, setting.Usage)
	}
	return nil
}

// Load parses the flags using arguments loaded by the argument loader.
func (s *source) Load(settings map[string]*congo.Setting) error {
	s.set.Parse(s.ArgLoader())
	return nil
}
