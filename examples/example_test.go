package examples

import (
	"fmt"

	"gitlab.com/silentteacup/congo"
	"gitlab.com/silentteacup/congo/sources/ini"
)

type Configuration struct {
	MaxUsers      uint  `name:"max-users" usage:"Sets how many users can be online concurrently."`
	Debug         bool  `name:"debug" usage:"En- or disables the debug mode"`
	InitialPoints int64 `name:"initial-points"`
	MaxPoints     uint64
}

func Example() {
	defaultCfg := Configuration{
		60,
		false,
		2000,
		3000000,
	}

	cfg := congo.New(
		"main",
		ini.FromFile("./important.ini"),
		ini.FromFile("./example.ini"),
	)
	cfg.Using(&defaultCfg)
	if err := cfg.Init(); err != nil {
		panic(err)
	}
	if err := cfg.Load(); err != nil {
		panic(err)
	}

	fmt.Printf("%+v", defaultCfg)
	// Output:
	// {MaxUsers:100 Debug:true InitialPoints:0 MaxPoints:3000000}
}
