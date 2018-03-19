package src_env

import (
	"congo"
	"fmt"
	"os"
)

func init() {
	os.Setenv("number", "54")
	os.Setenv("decimal", "0.5")
}

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
