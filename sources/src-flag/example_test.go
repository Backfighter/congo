package src_flag

import (
	"congo"
	"flag"
	"fmt"
	"os"
)

// Example is a basic example for the flag source
func Example() {
	// Set up arguments
	// This simulates arguments passed to the executable
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cmd", "-number=54", "-decimal=0.5", "args"}

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
