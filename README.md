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

	_, err2 := conftagz.Process(nil, &config)
	if err2 != nil {
		log.Fatalf("Config is bad: %v\n", err2)
	} else {
		fmt.Printf("Config good.\n")
	}

	fmt.Printf("Config: %v\n", config)

}
```


## What this is

An attempt to avoid repetitive, mundane, and slightly buggy code when reading and validating configuration files.

A common pattern with cloud apps is to specify a config file format in YAML or less commonly JSON. And then to:
- Parse that YAML into *Config* struct
- Check if all the values are valid
- Potentially override certain items in the struct with environmental variables if present
- Fill in defaults if values are completely missing

The order of operation and priority of these steps might differ here and there. But the gneral pattern is very common.

`conftagz` attempts to eleminate code writing for as much of this as possible by offering:

A `env:` struct tag which will replace the value of the field with the contents of the env var if present.

A `test:` struct tag which provides some basic tests (comprison, regex) or allows the calling of a custom func to check a value.

A `default:` struct tag which will replace any empty field with given value if no other method provides a value.




