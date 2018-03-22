package congo

import (
	"strconv"
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
//(Modified version of https://golang.org/src/flag/flag.go)

// -- bool Value

type boolValue bool

func newBoolValue(val bool, p *bool) Value {

	*p = val

	return (*boolValue)(p)

}

func (b *boolValue) Set(s string) error {

	v, err := strconv.ParseBool(s)

	*b = boolValue(v)

	return err

}

func (b *boolValue) Get() interface{} { return bool(*b) }

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) IsBoolFlag() bool { return true }

// -- int Value

type intValue int

func newIntValue(val int, p *int) Value {

	*p = val

	return (*intValue)(p)

}

func (i *intValue) Set(s string) error {

	v, err := strconv.ParseInt(s, 0, strconv.IntSize)

	*i = intValue(v)

	return err

}

func (i *intValue) Get() interface{} { return int(*i) }

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

// -- int64 Value

type int64Value int64

func newInt64Value(val int64, p *int64) Value {

	*p = val

	return (*int64Value)(p)

}

func (i *int64Value) Set(s string) error {

	v, err := strconv.ParseInt(s, 0, 64)

	*i = int64Value(v)

	return err

}

func (i *int64Value) Get() interface{} { return int64(*i) }

func (i *int64Value) String() string { return strconv.FormatInt(int64(*i), 10) }

// -- uint Value

type uintValue uint

func newUintValue(val uint, p *uint) Value {

	*p = val

	return (*uintValue)(p)

}

func (i *uintValue) Set(s string) error {

	v, err := strconv.ParseUint(s, 0, strconv.IntSize)

	*i = uintValue(v)

	return err

}

func (i *uintValue) Get() interface{} { return uint(*i) }

func (i *uintValue) String() string { return strconv.FormatUint(uint64(*i), 10) }

// -- uint64 Value

type uint64Value uint64

func newUint64Value(val uint64, p *uint64) Value {

	*p = val

	return (*uint64Value)(p)

}

func (i *uint64Value) Set(s string) error {

	v, err := strconv.ParseUint(s, 0, 64)

	*i = uint64Value(v)

	return err

}

func (i *uint64Value) Get() interface{} { return uint64(*i) }

func (i *uint64Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

// -- string Value

type stringValue string

func newStringValue(val string, p *string) Value {

	*p = val

	return (*stringValue)(p)

}

func (s *stringValue) Set(val string) error {

	*s = stringValue(val)

	return nil

}

func (s *stringValue) Get() interface{} { return string(*s) }

func (s *stringValue) String() string { return string(*s) }

// -- float64 Value

type float64Value float64

func newFloat64Value(val float64, p *float64) Value {

	*p = val

	return (*float64Value)(p)

}

func (f *float64Value) Set(s string) error {

	v, err := strconv.ParseFloat(s, 64)

	*f = float64Value(v)

	return err

}

func (f *float64Value) Get() interface{} { return float64(*f) }

func (f *float64Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }

// -- time.Duration Value

type durationValue time.Duration

func newDurationValue(val time.Duration, p *time.Duration) Value {

	*p = val

	return (*durationValue)(p)

}

func (d *durationValue) Set(s string) error {

	v, err := time.ParseDuration(s)

	*d = durationValue(v)

	return err

}

func (d *durationValue) Get() interface{} { return time.Duration(*d) }

func (d *durationValue) String() string { return (*time.Duration)(d).String() }

// Value is the interface to the dynamic value stored in a settings.
// (The default value is represented as a string.)
//
// Set is called whenever a setting is loaded from a source.
// The settings package may call the String method with a zero-valued receiver,
// such as a nil pointer.
type Value interface {
	String() string

	Set(string) error
}
