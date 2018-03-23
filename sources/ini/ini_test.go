package ini

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"errors"

	"gitlab.com/silentteacup/congo"
)

/*
Copyright (c) 2018 Peter Werner. All rights reserved.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

// TempFile represents a file that is temporally created for
// testing purposes.
type TempFile struct {
	name    string
	folder  string
	content string
}

// SimpleName gets the name of a file without the extensions.
func (t TempFile) SimpleName() string {
	return withoutExtension(t.name)
}

// Path gets the full path to the file.
func (t TempFile) Path() string {
	return filepath.Join(t.folder, t.name)
}

// withoutExtension removes the extension of a files path.
func withoutExtension(file string) string {
	return file[:strings.LastIndex(file, ".")]
}

// createTmpFiles writes the given temporary files to the disk and
// returns the paths to all written files.
func createTmpFiles(files ...TempFile) []string {
	created := make([]string, len(files))
	for i, f := range files {
		created[i] = f.Path()
		writeFile(f.Path(), []byte(f.content))
	}
	return created
}

// writeFile writes a files with given name and content.
// If it fails it panics.
func writeFile(file string, content []byte) {
	if err := ioutil.WriteFile(file, content, 0644); err != nil {
		panic(err)
	}
}

// cleanUp tries to delete given files from the disk.
// If it fails it panics.
func cleanUp(files ...string) {
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

var contentTmpl = TempFile{
	"test.ini",
	"./",
	content,
}

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

// TestFromFile tests the from file option for creating a ini source.
func TestFromFile(t *testing.T) {
	defer cleanUp(createTmpFiles(contentTmpl)...)
	v := &mockValue{
		nil,
		"",
		0,
		0,
	}
	settings := map[string]*congo.Setting{
		"number": {
			Name:     "number",
			Usage:    "usage",
			Value:    v,
			DefValue: "0",
		},
	}

	src := FromFile(contentTmpl.Path())
	if err := src.Init(settings); err != nil {
		t.Errorf("Expected to init without problems.\nBut got error: %s\n", err)
	}
	if v.SetCalls != 0 {
		t.Errorf("Expected Set() to not be called on mock.\nBut was called %d times.\n",
			v.SetCalls)
	}

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

type mockReaderCloser struct {
	*bytes.Buffer
}

func (*mockReaderCloser) Close() error {
	// Do nothing
	return nil
}

// TestNew tests the new option for creating a ini source.
func TestNew(t *testing.T) {
	v := &mockValue{
		nil,
		"",
		0,
		0,
	}
	settings := map[string]*congo.Setting{
		"number": {
			Name:     "number",
			Usage:    "usage",
			Value:    v,
			DefValue: "0",
		},
	}

	src := New(&mockReaderCloser{bytes.NewBufferString(content)})
	if err := src.Init(settings); err != nil {
		t.Errorf("Expected to init without problems.\nBut got error: %s\n", err)
	}
	if v.SetCalls != 0 {
		t.Errorf("Expected Set() to not be called on mock.\nBut was called %d times.\n",
			v.SetCalls)
	}

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

func TestIniSource_Load(t *testing.T) {
	s := FromFile("")
	err := s.Load(make(map[string]*congo.Setting))
	if err != nil {
		t.Errorf("Expected load with non-existent file to work without errors.\n"+
			"But returned error: %s\n", err)
	}
}

func TestIniSource_Load_NotLoose(t *testing.T) {
	s := FromFile("")
	err := s.SetLooseLoad(false).Load(make(map[string]*congo.Setting))
	if err == nil {
		t.Errorf("Expected non-loose load to return error.\n" +
			"But no error was returned.\n")
	}
}

func TestIniSource_Load_NoSection(t *testing.T) {
	s := FromBytes([]byte("")).Section("non-existent section")
	err := s.Load(make(map[string]*congo.Setting))
	if err != nil {
		t.Errorf("Expected load with non-existent section to work without errors.\n"+
			"But returned error: %s\n", err)
	}
}

func TestIniSource_Load_NoValue(t *testing.T) {
	v := &mockValue{
		nil,
		"",
		0,
		0,
	}
	v.Set("test")
	settings := map[string]*congo.Setting{
		"number": {
			Name:     "number",
			Usage:    "usage",
			Value:    v,
			DefValue: "0",
		},
	}
	s := FromBytes([]byte(""))
	err := s.Load(settings)
	if err != nil {
		t.Errorf("Expected load with non-existent value to work without errors.\n"+
			"But returned error: %s\n", err)
	}
}

func TestIniSource_Load_InvalidValue(t *testing.T) {
	v := &mockValue{
		errors.New("error"),
		"",
		0,
		0,
	}
	settings := map[string]*congo.Setting{
		"number": {
			Name:     "number",
			Usage:    "usage",
			Value:    v,
			DefValue: "0",
		},
	}
	s := FromBytes([]byte("number=invalid"))
	err := s.Load(settings)
	if err == nil {
		t.Errorf("Expected invalid value to cause return error on load.\n" +
			"But no error was returned.\n")
	}
}

// TestIniSource_WriteDefaults test the writing of defaults for
// the ini source.
func TestIniSource_WriteDefaults(t *testing.T) {
	v := &mockValue{
		nil,
		"",
		0,
		0,
	}
	v.Set("test")
	settings := map[string]*congo.Setting{
		"number": {
			Name:     "number",
			Usage:    "usage",
			Value:    v,
			DefValue: "0",
		},
	}
	s := FromBytes(make([]byte, 0))
	s.Init(settings)

	w := bytes.NewBufferString("")
	s.WriteDefaults(w)
	expected := "; usage\n" +
		"number = 0"
	actual := strings.Trim(w.String(), "\n ")
	if actual != expected {
		t.Errorf("Expected written default to be:\n %q\nBut was:\n %q\n",
			expected,
			actual)
	}
}
