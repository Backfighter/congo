package ini

import (
	"fmt"

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

// content is the content of a ini-file used in
// this example
const content = "" +
	"# Comment\n" +
	"; Number\n" +
	"number=54\n" +
	"; Decimal\n" +
	"decimal=0.5\n" +
	"[section]\n" +
	"duration=2h45m"

// Example is a basic example for the usage of the ini source.
func Example() {
	// Get ini source
	bytes := []byte(content)
	src := FromBytes(bytes)

	// main configuration
	cfg := congo.New("main", src)

	debug := cfg.Bool("debug", false, "Can be used to enable debug mode.")
	number := cfg.Int("number", 0, "Set a number")
	decimal := cfg.Float64("decimal", 0.2, "Set a decimal")

	// section of configuration
	subCfg := congo.New("section", src.Section("section"))
	duration := subCfg.Duration("duration", 0, "Set the duration.")

	// Load configurations
	cfg.Init()
	cfg.Load()
	subCfg.Init()
	subCfg.Load()

	if *debug {
		fmt.Println("Debug enabled!")
	}
	fmt.Printf("Using number %d and decimal %f\n", *number, *decimal)

	fmt.Printf("%v", *duration)
	//Output:
	//Using number 54 and decimal 0.500000
	//2h45m0s
}
