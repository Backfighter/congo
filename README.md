# Congo  
(*configure go*)
[![pipeline status](https://gitlab.com/SilentTeaCup/congo/badges/master/pipeline.svg)](https://gitlab.com/SilentTeaCup/congo/commits/master)
[![coverage report](https://gitlab.com/SilentTeaCup/congo/badges/master/coverage.svg)](https://gitlab.com/SilentTeaCup/congo/commits/master)

Congo is a configuration library based on the 
[flag package](https://golang.org/pkg/flag/) from the golang standard library.

It aims to be easy to use and extend.

## Show me the code

You can define your configuration using a struct.
```go
type Configuration struct {
	MaxUsers      uint  `name:"max-users" usage:"Sets how many users can be online concurrently."`
	Debug         bool  `name:"debug" usage:"En- or disables the debug mode"`
	InitialPoints int64 `name:"initial-points"`
	MaxPoints     uint64
}
```
Settings will be directly linked to your struct an will update
as soon as you Load() the configuration.
```go
func Example() {
	defaultCfg := Configuration{
		60,
		false,
		2000,
		3000000,
	}

	cfg := congo.New(
		"main", // name of the configuration
		ini.FromFile("./important.ini"), // sources added first will be preferred
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
```
Sources are prioritised in the oder they are passed to New().
Sources before others will overwrite the settings of the following sources.

### But I want none of this refection magic business!

No problem. Congo has you covered.
```go
func Example() {
	cfg := congo.New(
		"main", // name of the configuration
		ini.FromFile("./important.ini"), // sources added first will be preferred
		ini.FromFile("./example.ini"),
	)
	
    // Using assignment
    myInt := cfg.Int("my-int", 5, "Don't touch it, it's my int!")
    myDecimal := cfg.Float64("my-decimal", 0.8, "")
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
```
You can get a pointer to the setting yourself and manage them however you want.

## Sources

Congo uses modular sources to resolve settings. Currently the following
sources are supported:  
- INI files
- Environment variables
- Flags

But new ones can be added easily by implementing the source interface:  
```go
// Source is a source of settings e.g. flags, environment variables or a file.
type Source interface {
	// Init initializes this source fro given settings
	Init(map[string]*Setting) error
	// Load loads and sets the given settings
	Load(map[string]*Setting) error
}
```

Feel free to open MRs with new sources.

## Supported types

Congo supports the following types:
- int
- int64
- uint
- uint64 
- strings
- float64
- time.Duration
- Value

### Why is there no float/int/...32?

To avoid to much methods in the congo interface only allows the 64 bit versions since 
they carry the maximum amount of information an can be converted down if necessary.

### What is Value?

Value is an interfaces that can be used to add custom values for yourself:
```go
type Value interface {
	String() string

	Set(string) error
}
```
By implementing this interface you can add e.g. comma separated lists (a,b,c,d) or other 
custom values to the system.
```go
type Configuration struct {
	Custom     Value
}
//...
	defaultCfg := Configuration{
		&myCustomValue{[]int{1, 2, 3, 4}},
	}
//...
```
Custom can be set to any value you want and will be directly added as a setting like any other
field.