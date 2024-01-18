# conftagz

## Quick Example

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/IzumaNetworks/conftagz"
	"gopkg.in/yaml.v2"
)

type Config struct {
	WebhookURL string `yaml:"webhook_url" env:"APP_HOOK_URL" test:"~https://.*"`
	Port       int    `yaml:"port" env:"APP_PORT" default:"8888" test:">=1024,<65537"`
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


## What this is

An attempt to avoid repetitive, tedious, and slightly buggy code when reading and validating configuration files.

A common pattern with cloud apps is to specify a config file format in YAML, JSON or similar. And then to:
- Parse that YAML into some `MyConfig` struct which has struct tags for the parser
- Check if all the values are valid
- Potentially override certain items in the struct with environmental variables if present
- Set defaults on values with a *zero* value

The order of operation and priority of these steps might differ here and there. But the gneral pattern is common.

## The tags of `conftagz`

`conftagz` attempts to eleminate code writing for as much of this as possible by offering:

A `env:` struct tag which will replace the value of the field with the contents of the env var if present.

A `test:` struct tag which provides some basic tests (comprison, regex) or allows the calling of a custom func to check a value.

A `default:` struct tag which will replace any empty field with given value if no other method provides a value.

A `conf:` tag which can just change the behavior of `conftagz` itself for certain fields.

All tags are optional. Fields with no tag above are just ignored.

## Behavior and type support

`conftagz` behavior is specifcally designed to complement the behavior of the `yaml.v2` parser that almost everyone uses. 

Obviously, `conftagz` makes uses of the `reflect` package to do all this.

### Type support:

Fundamental types:
- `int`, `int16`, `int32`, `int64` and unsigned varaints
- `float32` and `float64`
- `string` ... `conftagz` uses the golang regex std library for regex tests
- pointers to all the above - `conftagz` will create the item if the pointer is nil _and_ a default or env var apply

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

`conftagz` will follow the struct pointer. If the struct is nil, it will create a new struct. This struct will have all zero values in it just as it were created with a `new()` call. This is necessary for go relfection to follow the struct and inspect its fields. If env vars stated are found it will assign their value to the field.

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

Sometimes a simple string value for a default won't cut it. Also, often defaults for structs and slices need more logic than a constant for an assignment. For this reason `default:` can call a registered `DefaultFunc`

```go
Field1        string        `yaml:"field1" default:"$(field1default)"`
```

and before calling `conftagz`:

```go
field1func := func(fieldname string) interface{} {
    return "field1funcval"
}

...

conftagz.RegisterDefaultFunc("field1default", field1func)
```

If `Field` is empty, then `field1func()` will be called, and it return value assigned.

## `test:` tag

The `test:` tag allows one or more tests to be performed on a field. By default, a call to `conftagz.Process()` will perform the tests _after_ all env vars and then defaults have been processed.

### Numeric fields

For numeric fields, `test:` supports: `>VAL`,`<VAL`,`>=VAL`,`<=VAL`,`==VAL`. Tests can be combined, comma seperated which will casue logical `&&` behavior.

For instance:

```go
	Port       int       `yaml:"port" test:">=1024,<65537"`
```

### String fields

String fields have regex support:

```go
	WebhookURL string    `yaml:"webhook_url" test:"~https://.*"`
```

Regex uses the standard regex golang library. The regex expression should start with a `~` to indicate its a regex expression. The expresison must `Regexp.Match()` the value or an error will be returned by `.Process()`

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

Custom functions allow various arbitrary tests. Becasue the function signature is the same regardless of type, the same function can be used for different types if needed.

## Running

The easiest way to use is:

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
- Runs the env subsituter: `EnvFieldSubstitution()`
- Runs the default subsiturer `SubsistuteDefaults()`
- Runs the tests `RunTestFlags()`

The order can be changed with the options.

Each of the above can also be called by itself. See test cases for more info.

## Examples

More docs to follow. See the `examples` folder for more examples.

Also refer to the `_test` files for more.
