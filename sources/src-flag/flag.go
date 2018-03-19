package src_flag

import (
	"congo"
	"flag"
	"os"
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
