package env

import (
	"fmt"
	"os"

	"gitlab.com/silentteacup/congo"
)

/*
Copyright (c) 2018 Peter Werner. All rights reserved.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

func init() {
	// Set up environment
	// This simulates a already setup environment
	os.Setenv("number", "54")
	os.Setenv("decimal", "0.5")
}

// Example is a basic example for the usage of the environment source.
func Example() {
	// Get environment source
	src := New()

	// Configuration
	cfg := congo.New("main", src)

	debug := cfg.Bool("debug", false, "Can be used to enable debug mode.")
	number := cfg.Int("number", 0, "Set a number")
	decimal := cfg.Float64("decimal", 0.2, "Set a decimal")

	// Load configurations
	cfg.Init()
	cfg.Load()

	if *debug {
		fmt.Println("Debug enabled!")
	}
	fmt.Printf("Using number %d and decimal %f\n", *number, *decimal)

	//Output:
	//Using number 54 and decimal 0.500000
}
