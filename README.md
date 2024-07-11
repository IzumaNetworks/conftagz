# conftagz

[![License](https://img.shields.io/:license-apache-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/go.izuma.io/conftagz)](https://goreportcard.com/report/go.izuma.io/IzumaNetworks/conftagz)
[![Build and Test](https://github.com/IzumaNetworks/conftagz/actions/workflows/build-and-test.yml/badge.svg)](https://github.com/IzumaNetworks/conftagz/actions/workflows/build-and-test.yml)

```
go get go.izuma.io/conftagz
```

## Quick Example

```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.izuma.io/conftagz"
	"gopkg.in/yaml.v2"
)

type Config struct {
	WebhookURL string `yaml:"webhook_url" env:"APP_HOOK_URL" test:"~https://.*"`
	Port       int    `yaml:"port" env:"APP_PORT" default:"8888" flag:"port" test:">=1024,<65537" usage:"Listen on port"`
	Expiration string `yaml:"expiration" default:"1h" test:"$(validtimeduration)"`
	DebugMode  bool   `yaml:"debug_mode" env:"DEBUG" flag:"debug"`
}

func ValidTimeDuration(val interface{}, fieldname string) bool {
	_, err := time.ParseDuration(val.(string))
	return err == nil
}

func main() {

	var config Config

	// load config file from yaml using yaml parser
	// Read the yaml file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Unmarshal the yaml file into the config struct
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// register that custom test
	conftagz.RegisterTestFunc("validtimeduration", ValidTimeDuration)

	// Run conftagz on the config struct
	// to validate the config, sub any env vars,
	// and put in defaults for missing items
	err2 := conftagz.Process(nil, &config)
	if err2 != nil {
		// some test tag failed
		log.Fatalf("Config is bad: %v\n", err2)
	} else {
		fmt.Printf("Config good.\n")
	}

	fmt.Printf("Config: %v\n", config)

}

```

Given a config file of:

```yaml
webhook_url: https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX
port: 8080
```

Will yield:
```
% ./example
Config good.
Config: {https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX 8080 1h false}
% ./example -debug
Config good.
Config: {https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX 8080 1h true}
% ./example -debug -port 8181
Config good.
Config: {https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX 8181 1h true}
% DEBUG=1 ./example
Config good.
Config: {https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX 8080 1h true}
% APP_PORT=8989 DEBUG=1 ./example
Config good.
Config: {https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX 8989 1h false}
% APP_PORT=89 ./example
2024/02/13 11:04:55 Config is bad: field Port: value 89 ! >= 1024
%  DEBUG=1 ./example --port 33
2024/02/13 11:05:07 Config is bad: field Port: value 33 ! >= 1024
```

## Motivation

There are many powerful and complicated libraries for configuration options and flags. [cobra](https://github.com/spf13/cobra), [viper](https://github.com/spf13/viper), [kong](https://github.com/alecthomas/kong), etc. But frankly software is already complex enough - and the last thing I wanted is some complicated library to just process command line arguments and config files. `conftagz` is the antithesis of these approaches.

When I go back to look at something from months ago - I want it to be super easy to figure out what's going on with the conf files and flags... Nor do I want to be confined to a specific way to layout components, or have to call dozens of library functions just to get the CLI options.

### Just use struct tags

Use structs + tags to define everything. Run `Process()` and that's it. No, it does not do even 1/8 the things cobra does. If you need that use cobra or one of the other fine options above.

A common pattern with cloud apps is to specify a config file format in YAML, JSON or similar - as a struct(s) in Go. And then to:
- Parse that YAML into some `MyConfig`-like struct
- Check if all the values are valid
- Override certain items with environmental variables if present
- Set defaults on values with a *zero* value
- Maybe replace some value with CLI flags

You can do all this with just struct tags using this package. Then make one call to `conftagz.Process()`

## The tags of `conftagz`

`conftagz` attempts to eleminate code writing for as much of this as possible by offering:

A `flag:` and (optional) `usage:` tag. This allows setting specific struct fields to be set by a command line flag. Uses the standrd `flag` package.

A `cflag:` `usage:` and `cobra:` flag allow you to use [cobra](https://github.com/spf13/cobra) instead of the normal stdlib `flags` package. See [Using Cobra for flags](#using-cobra-for-flags) section.

A `env:` struct tag which will replace the value of the field with the contents of the env var if present.

A `test:` struct tag which provides some basic validation (comparison, regex) or allows the calling of a custom func to check a value.

A `default:` struct tag which will replace any empty field with given value if no other method provides a value.

A `conf:` tag which can just change the behavior of `conftagz` itself for certain fields.

All tags are optional. Fields with no tag above are just ignored.

## Behavior and type support

`conftagz` behavior is specifically designed to complement the behavior of the `yaml.v2` parser that almost everyone uses.

Obviously, `conftagz` makes uses of the `reflect` package to do all this.

### Type support:

Fundamental types:
- `int`, `int16`, `int32`, `int64` and unsigned variants
- `bool` (not supported by `default:` tag as unnecessary)
- `float32` and `float64`
- `string` ... `conftagz` uses the golang regex std library for regex tests
- pointers to all the above - `conftagz` will create the item if the pointer is nil _and_ a default or env var are applied.

Structs & Slices
- Supports both and also their pointers
- Support for slices of structs and slices of pointers to structs 
- Default structs can be created if the yaml parser left a struct pointer nil by using a custom `DefaultFunc` like `default:"$(mydefaultfunc)"` See _custom defaults_
- `conftagz` will automatically create a new struct if the struct pointer is `nil`. This behavior can be avoided with `conf:"skip"` or `conf:"skipnil"`
- Nil slices of pointers to structs will be left alone without a custom function

Not supported
- Interfaces or `interface{}`
- `unintptr`
- Any other types not mentioned. Unsupported types are ignored.
- Anything which references itself. i.e. the config struct has a pointer pointing to itself

## `env:` tag

Example:

```go
	Port       int    `yaml:"port" env:"APP_PORT"`
```

`conftagz` will replace the field `Port` with the value of `APP_PORT` if the environmental variable `APP_PORT` exists. Type conversion from the string will happen automatically. If the env var is present _but_ it can not be converted for the type, an error is thrown. If the env var does not exist nothing will happen.

### Structs

Example:
```go
type Config struct {
    ...
    SSL        *SSLStuff `yaml:"sslstuff"`
    ...
}

type SSLStuff struct {
	Cert string `yaml:"cert" env:"SSL_CERT"`
	Key  string `yaml:"key" env:"SSL_KEY"`
}
```

`conftagz` will follow the struct pointer. If the struct is nil, it will create a new struct. This struct will have all zero values in it just as it were created with a `new()` call. This is necessary for go reflection to follow the struct and inspect its fields. If env vars stated are found it will assign their value to the field.

The behavior of automatically creating a struct from a `nil` pointer by the `conftagz` env substituter can be avoided with `skip`, `skipnil`, or `envskip` `conf:` tags:

```go
type Config struct {
    ...
    SSL        *SSLStuff `yaml:"sslstuff" conf:"envskip"`
    ...
}
```

## `default:` tag

The `default:` tag replaces _zero_ values of fields with `val` if a `default:"val"` tag exists. Type conversion takes place automatically just as with the `env:` tags. If the default tag is present _but_ it can not be converted for the type, an error is thrown. 

```go
type Config struct {
    Port       int       `yaml:"port" default:"8888"`
}
```

The `default:` tag is supported on fundamental types and slices of fundamentals:

```go
	SliceInts   []int   `yaml:"sliceints" default:"1,2,3"`
```

The above would fill an _empty_ SliceInts with `[1,2,3]`

As with the `env:` tags, the `conftagz` default substituter will follow the pointer. For fundamental types it will `new()` the given type and assign the default value to it - _if_ a default value is provided.

For struct pointers the substituter will _always_ create a `new()` struct unless told otherwise through a `conf:` `skip`, `skipnil`, or `defaultskip` tag.

Once the new struct is created, it will follow it and assign any defaults provided for each field.

### Default functions

Sometimes a simple string value for a default won't cut it. Also, often defaults for structs and slices need more logic than a constant for an assignment. For this reason `default:` can call a registered function meeting the `DefaultFunc` spec:

```go
Field1        string        `yaml:"field1" default:"$(field1default)"`
```

and before calling `conftagz` make sure the function is defined:

```go
field1func := func(fieldname string) interface{} {
    return "field1funcval"
}

Register it:
...

conftagz.RegisterDefaultFunc("field1default", field1func)
```

Then if `Field` is empty, then `field1func()` is called and its return value if assigned.

## `test:` tag

The `test:` tag allows one or more tests to be performed on a field. By default, a call to `conftagz.Process()` will perform the tests _after_ all env vars and then defaults have been processed.

### Numeric fields

For numeric fields, `test:` supports: `>VAL`,`<VAL`,`>=VAL`,`<=VAL`,`==VAL`. Tests can be combined, comma separated which will cause logical `&&` behavior.

For instance:

```go
	Port       int       `yaml:"port" test:">=1024,<65537"`
```

### String fields

String fields have regex support:

```go
	WebhookURL string    `yaml:"webhook_url" test:"~https://.*"`
```

Regex uses the standard regex golang library. The regex expression should start with a `~` to indicate its a regex expression. The expressions must `Regexp.Match()` the value or an error will be returned by `.Process()`

The regex is the only built-in test supported for string at the moment.

### Custom test functions

Like `default:`, `test:` support custom functions of the type `TestFunc` for tests on all supported types. For slices this is the only way to test.

Consider:

```go
type AStructWithCustom struct {
	Field1        string        `yaml:"field1" test:"$(field1test)"`
	DefaultStruct *InnerStruct2 `yaml:"inner" test:"$(fieldinnerstruct2test)"`
   	SliceInts     []int          `yaml:"sliceints" test:"$(sliceintstest)"`
}

field1func := func(val interface{}, fieldname string) bool {
    valstr, ok := val.(string)
    if !ok {
        // should never happen
        return false
    }
    if valstr != "stuff" {
        return false
    }
    return true
}

fieldstructfunc := func(val interface{}, fieldname string) bool {
    valstr, ok := val.(*InnerStruct2)
    if !ok {
        // should never happen
        return false
    }
    if valstr == nil || valstr.Stuff1 != "innerstuff" {
        return false
    }
    return true
}

testslicefunc := func(val interface{}, fieldname string) bool {
    valslice, ok := val.([]int)
    if !ok {
        t.Errorf("Expected slice, but got %v", val)
        return false
    }
    if len(valslice) < 3 {
        return false
    }
    if !(valslice[2] > valslice[1] && valslice[1] > valslice[0]) {
        return false
    }
    return true
}

RegisterTestFunc("sliceintstest", testslicefunc)
RegisterTestFunc("field1test", field1func)
RegisterTestFunc("fieldinnerstruct2test", fieldstructfunc)
```

Custom functions allow various arbitrary tests. Because the function signature is the same regardless of type, the same function can be used for different types if needed.

## Processing structs

The easiest way to use conftagz is:

```go
	err2 := conftagz.Process(nil, &config)
	if err2 != nil {
		// some test tag failed
        log.Fatalf("Config is bad: %v\n", err2)
	} else {
		fmt.Printf("Config good.\n")
	}
```

By default, `Process()` does the following in order:
- Runs the default subsiturer `SubsistuteDefaults()`
- Runs the env var subsituter: `EnvFieldSubstitution()`
- Runs the flag substiturer: `ProcessFlags()` or `PostProcessCobraFlags()` (if `PreProcessCobraFlags()` was called) 
- Runs the tests `RunTestFlags()`

The order can be changed with the options. By default command line switches if present override everything else.

Each of the above can also be called by itself. See test cases for more info.

## Using Cobra for flags

Given something like this:
```go
type Config struct {
	WebhookURL string    `yaml:"webhook_url" cflag:"webhookurl" usage:"URL to send webhooks to" cobra:"root"`
	Port       int       `yaml:"port" test:">=1024,<65537" cflag:"port" usage:"Port to listen on" cobra:"root"`
	SSL        *SSLStuff `yaml:"sslstuff"`
	Servers    []*Server `yaml:"servers"`
	// both a long and short --verbose or -v
	// cobra 'persistent' flag here
	Verbose bool `yaml:"verbose" env:"APP_VERBOSE" cflag:"verbose,v" usage:"Verbose output" cobra:"root,persistent"`
}

type AnotherStruct struct {
	AnotherField string `env:"ANOTHERFIELD" cflag:"anotherfield" cobra:"othercmd"`
}
```
Follow this general pattern to use cobra with conftagz:
```go
var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "A simple CLI application",
	RunE: func(cmd *cobra.Command, args []string) error {
		// implement your command
		...
		return nil
	},
}
// register your command with conftagz. Reference rootCmd with 'root' in your struct tag
conftagz.RegisterCobraCmd("root", rootCmd)
var otherCmd = &cobra.Command{
	Use:   "othercmd",
	Short: "Another command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
// another one
conftagz.RegisterCobraCmd("othercmd", otherCmd)

// run PreProcessCobraFlags for all struct with cobra tags
err = conftagz.PreProcessCobraFlags(&config, nil)
err = conftagz.PreProcessCobraFlags(&anotherstuct, nil)

rootCmd.AddCommand(otherCmd)
// Force cobra to parse the flags before running conftagz.Process
// You will need to parse all the flags for all the commands
// which have any conftagz fields
rootCmd.ParseFlags(os.Args)
otherCmd.ParseFlags(os.Args)

// Run conftagz on the structs
err2 := conftagz.Process(nil, &config)
err2 = conftagz.Process(nil, &anotherstuct)

// your structs should be filled in if flags were used
```

See `examples/examplecobra` for a fully working example.

## Examples

More docs to follow. See the `examples` folder for more examples.

Also refer to the test files for more.
