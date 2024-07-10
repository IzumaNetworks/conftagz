package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
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
	WebhookURL string    `yaml:"webhook_url" env:"APP_HOOK_URL" test:"~https://.*" cflag:"webhookurl" usage:"URL to send webhooks to" cobra:"root"`
	Port       int       `yaml:"port" env:"APP_PORT" default:"8888" test:">=1024,<65537" cflag:"port" usage:"Port to listen on" cobra:"root"`
	SSL        *SSLStuff `yaml:"sslstuff"`
	Servers    []*Server `yaml:"servers"`
	// LogSetup is a complex object
	// so we use the defaultLogSetup function to set it if its not already set
	// skipnil = means that if the pointer is nil we will not create what it can point to
	// nildefault = means that if the pointer is nil we will create what it can point to for 'detfault' tags
	// hence in this case a nil LogSetup will get whatever defaultLogSetup returns
	LogSetup *LogSetup `yaml:"log_setup" conf:"envskip,skipnil,nildefault" default:"$(defaultLogSetup)"`
	// for verbose we want to support both the --verbose and -v
	// verbose is a 'persistent' - see cobra docs for details. Basically this means the flag can be used on the 'root' command
	// and all its subcommands
	Verbose bool `yaml:"verbose" env:"APP_VERBOSE" cflag:"verbose,v" usage:"Verbose output" cobra:"root,persistent"`
}

type AnotherStruct struct {
	AnotherField string `env:"ANOTHERFIELD" cflag:"anotherfield" cobra:"othercmd"`
}

func RunMain() {
	var config Config

	var anotherstuct AnotherStruct

	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "A simple CLI application",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Running root command %+v\n", args)
			fmt.Printf("--------------------\n")
			fmt.Printf("Config: %+v\n", config)
			fmt.Printf("Logsetup: %+v\n", config.LogSetup)
			fmt.Printf("SSL: %+v\n", config.SSL)
			fmt.Printf("AnotherStruct: %+v\n", anotherstuct)
			fmt.Printf("--------------------\n")
			fmt.Printf("Verbose: %v\n", config.Verbose)
			return nil
		},
	}

	conftagz.RegisterCobraCmd("root", rootCmd)

	var otherCmd = &cobra.Command{
		Use:   "othercmd",
		Short: "Another command",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Running other command %+v\n", args)
			fmt.Printf("--------------------\n")
			fmt.Printf("Config: %+v\n", config)
			fmt.Printf("Logsetup: %+v\n", config.LogSetup)
			fmt.Printf("SSL: %+v\n", config.SSL)
			fmt.Printf("AnotherStruct: %+v\n", anotherstuct)
			fmt.Printf("--------------------\n")
			fmt.Printf("Verbose: %v\n", config.Verbose)
			return nil
		},
	}

	conftagz.RegisterCobraCmd("othercmd", otherCmd)

	// Root command flags
	// var boolFlag bool
	// var stringFlag string
	// var intFlag int
	// rootCmd.PersistentFlags().BoolVarP(&boolFlag, "bool", "b", false, "A boolean flag")
	// // rootCmd.PersistentFlags().StringVarP(&stringFlag, "string", "s", "", "A string flag")
	// // rootCmd.PersistentFlags().IntVarP(&intFlag, "int", "i", 0, "An int flag")

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
		fmt.Printf("called defaultLogSetupFunc\n")
		return &LogSetup{
			DebugPrefix:   "DEBUG",
			ErrorPrefix:   "ERROR",
			WarnPrefix:    "WARN",
			InfoPrefix:    "INFO",
			internalStuff: "internalthings",
		}
	}

	conftagz.RegisterDefaultFunc("defaultLogSetup", defaultLogSetupFunc)

	// PreProcessCobraFlags will add the flags to the cobra command
	// based on what is in the struct tags:
	err = conftagz.PreProcessCobraFlags(&config, nil)
	if err != nil {
		log.Fatalf("Unexpected error on PreProcessCobraFlags(config): %v", err)
	}
	// you can run this on any struct with cobra tags
	err = conftagz.PreProcessCobraFlags(&anotherstuct, nil)
	if err != nil {
		log.Fatalf("Unexpected error on PreProcessCobraFlags(anotherstuct): %v", err)
	}

	// make sure to add all commands before parsing flags
	rootCmd.AddCommand(otherCmd)
	// Force cobra to parse the flags before running conftagz.Process
	// You will need to parse all the flags for all the commands
	// which have any conftagz fields
	rootCmd.ParseFlags(os.Args)
	otherCmd.ParseFlags(os.Args)

	// Run conftagz on the config struct
	// to validate the config, sub any env vars, and put in defaults for missing items
	// pass in the optionn to use our own flag set
	// In the case of cobra tags - in options in os.Args will now be filled into the struct
	err2 := conftagz.Process(nil, &config)

	if err2 != nil {
		log.Fatalf("Config is bad: %v\n", err2)
	} else {
		fmt.Printf("Config good.\n")
	}

	// You can call conftagz on multiple structs
	err2 = conftagz.Process(nil, &anotherstuct)
	if err2 != nil {
		log.Fatalf("AnotherStruct is bad: %v\n", err2)
	} else {
		fmt.Printf("AnotherStruct good.\n")
	}

	fmt.Printf("AnotherField: %s\n", anotherstuct.AnotherField)

	rootCmd.Execute()
}

// RunMain is the main entry point for the application
// We make this app testable this way
func main() {
	RunMain()
}
