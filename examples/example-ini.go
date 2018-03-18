package main

import (
	"congo"
	"congo/sources/src-ini"
	"fmt"
)

func main() {
	cfg := congo.New("main", src_ini.New("./test.main"))
	debug := cfg.Bool("debug", false, "Can be used to enable debug mode.")
	number := cfg.Int("number", 0, "Set a number")
	decimal := cfg.Float64("decimal", 0.2, "Set a decimal")
	if *debug {
		fmt.Println("Debug enabled!")
	}
	fmt.Printf("Using number %d and decimal %f", *number, *decimal)
}
