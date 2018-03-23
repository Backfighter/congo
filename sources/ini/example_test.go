package ini

import (
	"fmt"

	"gitlab.com/silentteacup/congo"
)

/*
Copyright (c) 2018 Peter Werner. All rights reserved.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
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
