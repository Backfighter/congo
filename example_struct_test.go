package congo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
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

// customValue is a custom implementation of the Value interface.
// You can do your own custom implementation of the Value interface
// to add your own types that are not yet supported by congo.
type customValue struct {
	first  int
	second int
}

// String returns the string representation of the value.
// It must be usable by Set().
func (c *customValue) String() string {
	return fmt.Sprintf("%d:%d", c.first, c.second)
}

// Set sets the value accordingly given its string representation.
func (c *customValue) Set(in string) error {
	split := strings.Split(in, ":")
	if len(split) != 2 {
		return errors.New("value must consist of two parts separated by a ':'")
	}
	first, err := strconv.ParseInt(split[0], 0, 64)
	if err != nil {
		return err
	}
	second, err := strconv.ParseInt(split[1], 0, 64)
	if err != nil {
		return err
	}
	c.first = int(first)
	c.second = int(second)
	return nil
}

// configuration is the struct that will be used to generate
// the configuration.
type Configuration struct {
	// UpdateInterval is an setting for which name and usage are defined
	UpdateInterval time.Duration `name:"update-interval" usage:"Controls the time between updates"`
	// ExecutionPath is a setting for which only the name is defined
	// the usage message will be empty
	ExecutionPath string `name:"execution-path"`
	// MagicNumber is a setting that neither has a name nor a usage description.
	// The name will be directly derived from the field name ("MagicNumber").
	MagicNumber int
	// MagicDecimal is a settings that uses float64 as type.
	// Only float64 is supported since it can be converted to any
	// lower representation without loss.
	MagicDecimal float64
	// Custom is a custom value that represents a yet unknown type to congo
	Custom Value
}

// ExampleSource a basic implementation of the Source interface.
// Congo provides sources for: ini-files, environment variables and flags.
type ExampleSource struct{}

func (s *ExampleSource) Init(map[string]*Setting) error {
	// Do nothing (This source doesn't require initializing)
	return nil
}

func (s *ExampleSource) Load(settings map[string]*Setting) error {
	// This example source always provides the same values.
	// These could also be loaded from a file or other possible sources.
	//
	// Most sources will iterate over the settings and resolve them using their key.
	// After that they will set the value accordingly.

	if err := settings["update-interval"].Value.Set("1h"); err != nil {
		return err
	}
	if err := settings["MagicNumber"].Value.Set("0"); err != nil {
		return err
	}
	if err := settings["Custom"].Value.Set("9:8"); err != nil {
		return err
	}
	return nil
}

func Example() {
	defaultCfg := Configuration{
		UpdateInterval: time.Minute * 5,
		ExecutionPath:  "/execution/path",
		MagicNumber:    5,
		MagicDecimal:   0.8,
		Custom:         &customValue{2, 6},
	}

	cfg := New("main", &ExampleSource{})
	cfg.Using(&defaultCfg)
	if err := cfg.Init(); err != nil {
		panic(err)
	}
	if err := cfg.Load(); err != nil {
		panic(err)
	}

	fmt.Printf("%+v", defaultCfg)
	// Output:
	// {UpdateInterval:1h0m0s ExecutionPath:/execution/path MagicNumber:0 MagicDecimal:0.8 Custom:9:8}
}
