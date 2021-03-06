package congo

import (
	"bytes"
	"testing"
	"time"
)

/*
Copyright (c) 2018 Peter Werner. All rights reserved.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
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

func newMockValue(err error) *mockValue {
	return &mockValue{
		err,
		"",
		0,
		0,
	}
}

func setupTestCongo() (Congo, *testSource) {
	s := &testSource{}
	sources := []Source{s}
	output := bytes.NewBufferString("")
	c := congo{
		sources,
		make(map[string]*Setting),
		"test",
		output,
	}
	return &c, s
}

func TestCongo_Bool(t *testing.T) {
	c, s := setupTestCongo()
	v := c.Bool("test", false, "Usage")
	c.Load()
	param, ok := s.LoadParam["test"]
	if !ok {
		t.Errorf("Expected load parameters to contain %s.\nBut didn't extist.\n", "test")
		return
	}
	param.Value.Set("true")
	if *v != true {
		t.Errorf("Expected value to be set to true.\nBut was set to %v.\n", *v)
	}
}

func TestCongo_Int64(t *testing.T) {
	c, s := setupTestCongo()
	v := c.Int64("test", -23429284, "Usage")
	c.Load()
	param, ok := s.LoadParam["test"]
	if !ok {
		t.Errorf("Expected load parameters to contain %s.\nBut didn't extist.\n", "test")
		return
	}
	param.Value.Set("5")
	if *v != 5 {
		t.Errorf("Expected value to be set to %d.\nBut was set to %d.\n", 5, *v)
	}
}

func TestCongo_Uint(t *testing.T) {
	c, s := setupTestCongo()
	v := c.Uint("test", 72938479284, "Usage")
	c.Load()
	param, ok := s.LoadParam["test"]
	if !ok {
		t.Errorf("Expected load parameters to contain %s.\nBut didn't extist.\n", "test")
		return
	}
	param.Value.Set("9")
	if *v != 9 {
		t.Errorf("Expected value to be set to %d.\nBut was set to %d.\n", 9, *v)
	}
}

func TestCongo_Uint64(t *testing.T) {
	c, s := setupTestCongo()
	v := c.Uint64("test", 0, "Usage")
	c.Load()
	param, ok := s.LoadParam["test"]
	if !ok {
		t.Errorf("Expected load parameters to contain %s.\nBut didn't extist.\n", "test")
		return
	}
	param.Value.Set("8980980")
	if *v != 8980980 {
		t.Errorf("Expected value to be set to %d.\nBut was set to %d.\n", 8980980, *v)
	}
}

func TestCongo_String(t *testing.T) {
	c, s := setupTestCongo()
	v := c.String("test", "string", "Usage")
	c.Load()
	param, ok := s.LoadParam["test"]
	if !ok {
		t.Errorf("Expected load parameters to contain %s.\nBut didn't extist.\n", "test")
		return
	}
	param.Value.Set("new string")
	if *v != "new string" {
		t.Errorf("Expected value to be set to %s.\nBut was set to %s.\n", "new string", *v)
	}
}

func TestCongo_Duration(t *testing.T) {
	c, s := setupTestCongo()
	v := c.Duration("test", time.Hour*3, "Usage")
	c.Load()
	param, ok := s.LoadParam["test"]
	if !ok {
		t.Errorf("Expected load parameters to contain %s.\nBut didn't extist.\n", "test")
		return
	}
	param.Value.Set("5h3m")
	if *v != time.Hour*5+time.Minute*3 {
		t.Errorf("Expected value to be set to %v.\nBut was set to %v.\n",
			time.Hour*5+time.Minute*3, *v)
	}
}

func TestCongo_Var(t *testing.T) {
	defer testForPanic(t)
	c, _ := setupTestCongo()
	v := newMockValue(nil)
	c.Var(v, "test", "usage")
	c.Var(v, "test", "usage")
}

type testStruct struct {
	VBool    bool  `name:"value-bool" usage:"bool usage"`
	VInt     int   `name:"int-bool"`
	VInt64   int64 `usage:"int64 usage"`
	VUint    uint
	VUint64  uint64
	VString  string
	VFloat64 float64
	time.Duration
	Value
	Vnil Value
	// Unknown type
	float32
	// Unexported
	bool
}

func TestCongo_Using(t *testing.T) {
	v := newMockValue(nil)
	v.SetParam = "test"
	settings := testStruct{VUint: 5, Value: v}
	expected := []struct {
		name     string
		usage    string
		defValue string
	}{
		{"value-bool", "bool usage", "false"},
		{"int-bool", "", "0"},
		{"VInt64", "int64 usage", "0"},
		{"VUint", "", "5"},
		{"VUint64", "", "0"},
		{"VString", "", ""},
		{"VFloat64", "", "0"},
		{"Duration", "", "0s"},
		{"Value", "", "test"},
	}
	c, s := setupTestCongo()
	c.Using(&settings)
	c.Init()
	params := s.InitParam
	if len(params) != len(expected) {
		t.Errorf("To much/less settings interpreted. Expected: %d\n"+
			"But got: %d\n", len(expected), len(params))
	}
	for _, e := range expected {
		s, ok := params[e.name]
		if !ok {
			t.Errorf("Expected %q to be in the settings.\nBut was not.\n", e.name)
			continue
		}
		if s.Name != e.name {
			t.Errorf("Expected name to be: %s\nBut got: %s\n", e.name, s.Name)
		}
		if s.DefValue != e.defValue {
			t.Errorf("Expected default value for %q to be: %s\n"+
				"But got: %s\n", e.name, e.defValue, s.DefValue)
		}
		if s.Usage != e.usage {
			t.Errorf("Expected usage for %q to be: %s\n"+
				"But got: %s\n", e.name, e.usage, s.Usage)
		}
	}
}

func TestCongo_Using2(t *testing.T) {
	defer testForPanic(t)
	c, _ := setupTestCongo()
	c.Using(5)
}

func testForPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("Expected code to panic but it didn't.")
	}
}
