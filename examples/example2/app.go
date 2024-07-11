package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go.izuma.io/conftagz"
	"gopkg.in/yaml.v2"
)

// Example with multiple structs

type Server struct {
	IP   string `yaml:"ip" test:"~[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$"`
	Name string `yaml:"name" test:"~.+"`
}
type SSLStuff struct {
	// cert and key should be at least 10 chars long
	Cert string `yaml:"cert" env:"SSL_CERT" test:"~.{10,}"`
	Key  string `yaml:"key" env:"SSL_KEY" test:"~.{10,}"`
}

type LogSetup struct {
	DebugPrefix   string `yaml:"debug_prefix"`
	ErrorPrefix   string `yaml:"error_prefix"`
	WarnPrefix    string `yaml:"warn_prefix"`
	InfoPrefix    string `yaml:"info_prefix"`
	internalStuff string
}

type Config struct {
	WebhookURL string    `yaml:"webhook_url" env:"APP_HOOK_URL" test:"~https://.*" flag:"webhookurl" usage:"URL to send webhooks to"`
	Port       int       `yaml:"port" env:"APP_PORT" default:"8888" test:">=1024,<65537"`
	SSL        *SSLStuff `yaml:"sslstuff"`
	Servers    []*Server `yaml:"servers"`
	LogSetup   *LogSetup `yaml:"log_setup" conf:"envskip" default:"$(defaultLogSetup)"`
}

type AnotherStruct struct {
	AnotherField string `env:"ANOTHERFIELD" flag:"anotherfield"`
}

// Make the app testable this way
func RunMain() {
	var config Config

	flagset := flag.NewFlagSet("test", flag.ContinueOnError)
	someotherflag := flagset.Bool("someotherflag", false, "Test if a flag not in confg struct is fine")
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

	defaultLogSetupFunc := func(fieldname string) interface{} {
		return &LogSetup{
			DebugPrefix:   "DEBUG",
			ErrorPrefix:   "ERROR",
			WarnPrefix:    "WARN",
			InfoPrefix:    "INFO",
			internalStuff: "internalthings",
		}
	}

	conftagz.RegisterDefaultFunc("defaultLogSetup", defaultLogSetupFunc)

	var anotherstruct AnotherStruct

	conftagz.PreProcessFlagsWithFlagSet(&anotherstruct, flagset)

	// Run conftagz on the config struct
	// to validate the config, sub any env vars, and put in defaults for missing items
	// pass in the optionn to use our own flag set
	err2 := conftagz.Process(&conftagz.ConfTagOpts{
		FlagTagOpts: &conftagz.FlagFieldSubstOpts{
			UseFlags: flagset,
		},
	}, &config)

	if err2 != nil {
		log.Fatalf("Config is bad: %v\n", err2)
	} else {
		fmt.Printf("Config good.\n")
	}

	// You can call conftagz on multiple structs also
	// and with the same flag options if needed
	err2 = conftagz.Process(&conftagz.ConfTagOpts{
		FlagTagOpts: &conftagz.FlagFieldSubstOpts{
			UseFlags: flagset,
		},
	}, &anotherstruct)

	if err2 != nil {
		log.Fatalf("AnotherStruct is bad: %v\n", err2)
	} else {
		fmt.Printf("AnotherStruct good.\n")
	}

	if *someotherflag {
		fmt.Printf("someotherflag is set\n")
	}

	if anotherstruct.AnotherField != "" {
		fmt.Printf("AnotherField: %v\n", anotherstruct.AnotherField)
	}

	fmt.Printf("Config: %+v\n", config)
	fmt.Printf("Logsetup: %+v\n", config.LogSetup)
	fmt.Printf("SSL: %+v\n", config.SSL)

}

func main() {
	RunMain()
}
