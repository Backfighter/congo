package src_ini

import (
	"congo"
	"fmt"
)

const content = "" +
	"# Comment\n" +
	"; Number\n" +
	"number=54\n" +
	"; Decimal\n" +
	"decimal=0.5\n" +
	"[section]\n" +
	"duration=2h45m"

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

	// load configurations
	cfg.Load()
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
