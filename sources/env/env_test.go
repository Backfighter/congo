package env

import (
	"os"
	"strings"
	"testing"

	"gitlab.com/silentteacup/congo"
)

/*
Copyright (c) 2018 Peter Werner. All rights reserved.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
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
