package env

import (
	"os"
	"strings"
	"testing"

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

type mockValue struct {
	returnErr   error
	SetParam    string
	StringCalls int
	SetCalls    int
}

func (v *mockValue) String() string {
	v.StringCalls++
	return v.SetParam
}

func (v *mockValue) Set(param string) error {
	v.SetCalls++
	v.SetParam = param
	return v.returnErr
}

func (v *mockValue) Reset() {
	v.StringCalls = 0
	v.SetCalls = 0
}

func TestSource_WithTranslator(t *testing.T) {
	src := New().WithTranslator(func(s string) []string {
		return []string{s, strings.ToLower(s)}
	})
	v := &mockValue{}
	settings := map[string]*congo.Setting{
		"NUMBER": {
			"NUMBER",
			"",
			v,
			"0",
		},
	}
	if err := src.Init(settings); err != nil {
		t.Errorf("Expected to init without problems.\nBut got error: %s\n", err)
	}

	// Set up environment
	// This simulates a already setup environment
	os.Setenv("number", "54")
	os.Setenv("NUMBER", "5")
	if err := src.Load(settings); err != nil {
		t.Errorf("Expected to load without problems.\nBut got error: %s\n", err)
	}
	if v.SetCalls != 1 {
		t.Errorf("Expected Set() to be called once on mock.\nBut was called %d times.\n",
			v.SetCalls)
	}
	if v.SetParam != "5" {
		t.Errorf("Expected Set() to be called with %s as parameter"+
			".\nBut was called with %s.\n", "5", v.SetParam)
	}

	os.Unsetenv("NUMBER")
	v.Reset()
	if err := src.Load(settings); err != nil {
		t.Errorf("Expected to load without problems.\nBut got error: %s\n", err)
	}
	if v.SetCalls != 1 {
		t.Errorf("Expected Set() to be called once on mock.\nBut was called %d times.\n",
			v.SetCalls)
	}
	if v.SetParam != "54" {
		t.Errorf("Expected Set() to be called with %s as parameter"+
			".\nBut was called with %s.\n", "54", v.SetParam)
	}
}

func TestPrefixSdtTranslator(t *testing.T) {
	translator := PrefixSdtTranslator("pref@x")

	result := translator(" strinG124öä-;#")

	expected := "PREFX_STRING124_"

	if result[0] != expected {
		t.Errorf("Expected translation to be %q\nBut got: %q\n", expected, result[0])
	}
}
