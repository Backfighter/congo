package src_ini

import (
	"bytes"
	"congo"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
