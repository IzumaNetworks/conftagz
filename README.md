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
	_, err2 := conftagz.Process(nil, &config)
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

More docs to follow. See the `examples` folder for more examples.
