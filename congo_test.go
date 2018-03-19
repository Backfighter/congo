package congo

import (
	"bytes"
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

type testSource struct {
	InitParam map[string]*Setting
	LoadParam map[string]*Setting
	InitErr   error
	LoadErr   error
}

func (t *testSource) Init(param map[string]*Setting) error {
	t.InitParam = param
	return t.InitErr
}

func (t *testSource) Load(param map[string]*Setting) error {
	t.LoadParam = param
	return t.LoadErr
}

func TestCongo_Bool(t *testing.T) {
	source := &testSource{}
	sources := []Source{source}
	output := bytes.NewBufferString("")
	c := congo{
		sources,
		make(map[string]*Setting),
		"test",
		output,
	}
	v := c.Bool("test", false, "Usage")
	c.Load()
	param, ok := source.LoadParam["test"]
	if !ok {
		t.Errorf("Expected load parameters to contain %s.\nBut didn't extist.\n", "test")
		return
	}
	param.Value.Set("true")
	if *v != true {
		t.Errorf("Expected value to be set to true.\nBut was set to %v.\n", *v)
	}
}
