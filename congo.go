package congo

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"time"
)

/*
Copyright (c) 2009 The Go Authors. All rights reserved.
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
//(Modified version of https://golang.org/src/setting/setting.go)

// Source is a source of settings e.g. flags, environment variables or a file.
type Source interface {
	// Init initializes this source fro given settings
	Init(map[string]*Setting) error
	// Load loads and sets the given settings
	Load(map[string]*Setting) error
}

// Setting is a value that is part of the configuration of a system.
type Setting struct {
	Name     string // name the key of the setting
	Usage    string // contains information on how to use the setting
	Value    Value  // value as set
	DefValue string // default value (as text)
}

// New creates a new configuration that uses given sources to resolve
// the settings. Sources will be prioritized based on the order they are
// given to this function. Sources in the front will overwrite settings
// provides by sources in the back.
func New(name string, sources ...Source) Congo {
	return &congo{
		sources,
		make(map[string]*Setting),
		name,
		os.Stderr,
	}
}

// Congo is a configuration capable of loading settings from different
// sources.
type Congo interface {
	// BoolVar defines a bool setting with specified name, default value, and usage string.
	// The argument p points to a bool variable in which to store the value of the setting.
	BoolVar(p *bool, name string, value bool, usage string)
	// Bool defines a bool setting with specified name, default value, and usage string.
	// The return value is the address of a bool variable that stores the value of the setting.
	Bool(name string, value bool, usage string) *bool

	// IntVar defines an int setting with specified name, default value, and usage string.
	// The argument p points to an int variable in which to store the value of the setting.
	IntVar(p *int, name string, value int, usage string)
	// Int defines an int setting with specified name, default value, and usage string.
	// The return value is the address of an int variable that stores the value of the setting.
	Int(name string, value int, usage string) *int

	// Int64Var defines an int64 setting with specified name, default value, and usage string.
	// The argument p points to an int64 variable in which to store the value of the setting.
	Int64Var(p *int64, name string, value int64, usage string)
	// Int64 defines an int64 setting with specified name, default value, and usage string.
	// The return value is the address of an int64 variable that stores the value of the setting.
	Int64(name string, value int64, usage string) *int64

	// UintVar defines a uint setting with specified name, default value, and usage string.
	// The argument p points to a uint variable in which to store the value of the setting.
	UintVar(p *uint, name string, value uint, usage string)
	// Uint64 defines a uint64 setting with specified name, default value, and usage string.
	// The return value is the address of a uint64 variable that stores the value of the setting.
	Uint64(name string, value uint64, usage string) *uint64

	// StringVar defines a string setting with specified name, default value, and usage string.
	// The argument p points to a string variable in which to store the value of the setting.
	StringVar(p *string, name string, value string, usage string)
	// String defines a string setting with specified name, default value, and usage string.
	// The return value is the address of a string variable that stores the value of the setting.
	String(name string, value string, usage string) *string

	// Float64Var defines a float64 setting with specified name, default value, and usage string.
	// The argument p points to a float64 variable in which to store the value of the setting.
	Float64Var(p *float64, name string, value float64, usage string)
	// Float64 defines a float64 setting with specified name, default value, and usage string.
	// The return value is the address of a float64 variable that stores the value of the setting.
	Float64(name string, value float64, usage string) *float64

	// DurationVar defines a time.Duration setting with specified name, default value, and usage string.
	// The argument p points to a time.Duration variable in which to store the value of the setting.
	// The setting accepts a value acceptable to time.ParseDuration.
	DurationVar(p *time.Duration, name string, value time.Duration, usage string)
	// Duration defines a time.Duration setting with specified name, default value, and usage string.
	// The return value is the address of a time.Duration variable that stores the value of the setting.
	// The setting accepts a value acceptable to time.ParseDuration.
	Duration(name string, value time.Duration, usage string) *time.Duration

	// Var defines a setting with the specified name and usage string. The type and
	// value of the setting are represented by the first argument, of type Value, which
	// typically holds a user-defined implementation of Value. For instance, the
	// caller could create a setting that turns a comma-separated string into a slice
	// of strings by giving the slice the methods of Value; in particular, Set would
	// decompose the comma-separated string into the slice.
	Var(value Value, name string, usage string)

	// Init initializes the configuration sources.
	Init() error

	// Load loads the configuration from the sources.
	Load() error
}

type congo struct {
	sources  []Source            // sources for the settings
	settings map[string]*Setting // settings
	name     string              // name of the configuration
	output   io.Writer
}

// BoolVar defines a bool setting with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the setting.
func (c *congo) BoolVar(p *bool, name string, value bool, usage string) {
	c.Var(NewBoolValue(value, p), name, usage)
}

// Bool defines a bool setting with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the setting.
func (c *congo) Bool(name string, value bool, usage string) *bool {
	p := new(bool)
	c.BoolVar(p, name, value, usage)
	return p
}

// IntVar defines an int setting with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the setting.
func (c *congo) IntVar(p *int, name string, value int, usage string) {
	c.Var(NewIntValue(value, p), name, usage)
}

// Int defines an int setting with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the setting.
func (c *congo) Int(name string, value int, usage string) *int {
	p := new(int)
	c.IntVar(p, name, value, usage)
	return p
}

// Int64Var defines an int64 setting with specified name, default value, and usage string.
// The argument p points to an int64 variable in which to store the value of the setting.
func (c *congo) Int64Var(p *int64, name string, value int64, usage string) {
	c.Var(NewInt64Value(value, p), name, usage)
}

// Int64 defines an int64 setting with specified name, default value, and usage string.
// The return value is the address of an int64 variable that stores the value of the setting.
func (c *congo) Int64(name string, value int64, usage string) *int64 {
	p := new(int64)
	c.Int64Var(p, name, value, usage)
	return p
}

// UintVar defines a uint setting with specified name, default value, and usage string.
// The argument p points to a uint variable in which to store the value of the setting.
func (c *congo) UintVar(p *uint, name string, value uint, usage string) {
	c.Var(NewUintValue(value, p), name, usage)
}

// Uint defines a uint setting with specified name, default value, and usage string.
// The return value is the address of a uint variable that stores the value of the setting.
func (c *congo) Uint(name string, value uint, usage string) *uint {
	p := new(uint)
	c.UintVar(p, name, value, usage)
	return p
}

// Uint64Var defines a uint64 setting with specified name, default value, and usage string.
// The argument p points to a uint64 variable in which to store the value of the setting.
func (c *congo) Uint64Var(p *uint64, name string, value uint64, usage string) {
	c.Var(NewUint64Value(value, p), name, usage)
}

// Uint64 defines a uint64 setting with specified name, default value, and usage string.
// The return value is the address of a uint64 variable that stores the value of the setting.
func (c *congo) Uint64(name string, value uint64, usage string) *uint64 {
	p := new(uint64)
	c.Uint64Var(p, name, value, usage)
	return p
}

// StringVar defines a string setting with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the setting.
func (c *congo) StringVar(p *string, name string, value string, usage string) {
	c.Var(NewStringValue(value, p), name, usage)
}

// String defines a string setting with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the setting.
func (c *congo) String(name string, value string, usage string) *string {
	p := new(string)
	c.StringVar(p, name, value, usage)
	return p
}

// Float64Var defines a float64 setting with specified name, default value, and usage string.
// The argument p points to a float64 variable in which to store the value of the setting.
func (c *congo) Float64Var(p *float64, name string, value float64, usage string) {
	c.Var(NewFloat64Value(value, p), name, usage)
}

// Float64 defines a float64 setting with specified name, default value, and usage string.
// The return value is the address of a float64 variable that stores the value of the setting.
func (c *congo) Float64(name string, value float64, usage string) *float64 {
	p := new(float64)
	c.Float64Var(p, name, value, usage)
	return p
}

// DurationVar defines a time.Duration setting with specified name, default value, and usage string.
// The argument p points to a time.Duration variable in which to store the value of the setting.
// The setting accepts a value acceptable to time.ParseDuration.
func (c *congo) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	c.Var(NewDurationValue(value, p), name, usage)
}

// Duration defines a time.Duration setting with specified name, default value, and usage string.
// The return value is the address of a time.Duration variable that stores the value of the setting.
// The setting accepts a value acceptable to time.ParseDuration.
func (c *congo) Duration(name string, value time.Duration, usage string) *time.Duration {
	p := new(time.Duration)
	c.DurationVar(p, name, value, usage)
	return p
}

// Var defines a setting with the specified name and usage string. The type and
// value of the setting are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a setting that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func (c *congo) Var(value Value, name string, usage string) {
	// Remember the default value as a string; it won't change.
	setting := &Setting{name, usage, value, value.String()}
	_, alreadythere := c.settings[name]
	if alreadythere {
		var msg string
		if c.name == "" {
			msg = fmt.Sprintf("setting redefined: %s", name)
		} else {
			msg = fmt.Sprintf("%s setting redefined: %s", c.name, name)
		}
		fmt.Fprintln(c.output, msg)
		panic(msg) // Happens only if settings are declared with identical names
	}
	if c.settings == nil {
		c.settings = make(map[string]*Setting)
	}
	c.settings[name] = setting
}

// Init initializes the configuration sources.
func (c *congo) Init() error {
	for i := len(c.sources) - 1; i >= 0; i-- {
		if err := c.sources[i].Init(c.settings); err != nil {
			return err
		}
	}
	return nil
}

// Load loads the configuration from the sources.
func (c *congo) Load() error {
	for i := len(c.sources) - 1; i >= 0; i-- {
		if err := c.sources[i].Load(c.settings); err != nil {
			return err
		}
	}
	return nil
}

// Using takes an arbitrary struct and turns it into a configuration.
// Fields of the struct are read and linked to the configuration.
// Values of the fields are updated as soon as Load() is called.
// Values already assigned to a field wil be used as default value for the setting.
//
// Fields can be annotated with tags describing the setting. Available tags are:
//
// `name`: Will be used as name. If not present the name of the field will be used.
//
// `usage`: Will be used as usage message (can be omitted).
//
// Supported types for field are: int, int64, uint, uint64, strings, float64, time.Duration
// and Value.
// A field that implements the Value type can be used to add custom, yet unsupported types.
// These fields will be directly added using the Var() method.
//
// All other types will be ignored!
func (c *congo) Using(configurationStruct interface{}) {
	v := reflect.ValueOf(configurationStruct)
	for i := 0; i < v.NumField(); i++ {
		c.register(v.Type().Field(i), v.Field(i).Interface())
	}
}

const (
	usageTag = "usage"
	nameTag  = "name"
)

// register registers a StructField with given value into the settings
// the type of the value is converted into a Value and added as settings
// using additional information from tags.
func (c *congo) register(f reflect.StructField, v interface{}) {
	usage := f.Tag.Get(usageTag)
	name, ok := f.Tag.Lookup(nameTag)
	if !ok {
		name = f.Name
	}
	switch a := v.(type) {
	case bool:
		c.BoolVar(&a, name, a, usage)
	case int:
		c.IntVar(&a, name, a, usage)
	case int64:
		c.Int64Var(&a, name, a, usage)
	case uint:
		c.UintVar(&a, name, a, usage)
	case uint64:
		c.Uint64Var(&a, name, a, usage)
	case string:
		c.StringVar(&a, name, a, usage)
	case float64:
		c.Float64Var(&a, name, a, usage)
	case time.Duration:
		c.DurationVar(&a, name, a, usage)
	case Value:
		c.Var(a, name, usage)
	default:
		// Do nothing.
	}
}
