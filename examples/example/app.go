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
	Port       int    `yaml:"port" env:"APP_PORT" default:"8888" flag:"port" test:">=1024,<65537" usage:"Port to listen on"`
	Expiration string `yaml:"expiration" default:"1h" test:"$(validtimeduration)"`
	DebugMode  bool   `yaml:"debug_mode" env:"DEBUG" flag:"debug"`
}

func ValidTimeDuration(val interface{}, fieldname string) bool {
	_, err := time.ParseDuration(val.(string))
	return err == nil
}

func RunMain() {

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

func main() {
	RunMain()
}
