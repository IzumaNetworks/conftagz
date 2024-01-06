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
		log.Fatalf("Error from conftags: %v\n", err2)
	}

	fmt.Printf("done. %v\n", config)

}
