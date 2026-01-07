package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.izuma.io/conftagz"
)

// Example demonstrating GetDefaultOnlyStruct
// This function creates a new struct with only default values applied
// It also demonstrates using custom default functions via $(funcName) syntax

// DatabaseConfig represents database connection settings
type DatabaseConfig struct {
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
	Username string `default:"postgres"`
	Password string `default:"secret"`
	Database string `default:"myapp"`
	SSLMode  string `default:"disable"`
}

// ServerConfig represents HTTP server settings
type ServerConfig struct {
	ListenAddr     string `default:"0.0.0.0"`
	Port           int    `default:"8080"`
	ReadTimeout    int    `default:"30"`
	WriteTimeout   int    `default:"30"`
	MaxHeaderBytes int    `default:"1048576"`
}

// LogConfig represents logging settings
type LogConfig struct {
	Level  string   `default:"info"`
	Format string   `default:"json"`
	Output []string `default:"stdout,file"`
}

// FeatureFlags represents feature toggles - populated by a custom function
type FeatureFlags struct {
	EnableMetrics    bool
	EnableProfiling  bool
	EnableTracing    bool
	MaintenanceMode  bool
	MaxConcurrency   int
	BuildInfo        string
}

// RuntimeInfo represents runtime information - populated by a custom function
type RuntimeInfo struct {
	Hostname  string
	StartTime string
	PID       int
}

// AppConfig is the main configuration struct with nested configs
type AppConfig struct {
	AppName  string          `default:"MyApplication"`
	Version  string          `default:"1.0.0"`
	Debug    bool            `default:"false"`
	Database *DatabaseConfig // pointer to nested struct
	Server   ServerConfig    // embedded struct
	Log      *LogConfig      // pointer to nested struct
	Features *FeatureFlags   `default:"$(defaultFeatureFlags)"` // custom function returns struct
	Runtime  *RuntimeInfo    `default:"$(defaultRuntimeInfo)"`  // another custom function
}

// defaultFeatureFlags is a custom default function that returns feature flags
// This demonstrates how to use dynamic defaults that can be computed at runtime
func defaultFeatureFlags(fieldname string) interface{} {
	return &FeatureFlags{
		EnableMetrics:   true,
		EnableProfiling: false,
		EnableTracing:   true,
		MaintenanceMode: false,
		MaxConcurrency:  100,
		BuildInfo:       "v1.0.0-abc123",
	}
}

// defaultRuntimeInfo is a custom default function that returns runtime info
// This shows how defaults can include dynamic values like hostname and time
func defaultRuntimeInfo(fieldname string) interface{} {
	hostname, _ := os.Hostname()
	return &RuntimeInfo{
		Hostname:  hostname,
		StartTime: time.Now().Format(time.RFC3339),
		PID:       os.Getpid(),
	}
}

// RunMain demonstrates the GetDefaultOnlyStruct functionality
func RunMain() {
	// Register custom default functions before calling GetDefaultOnlyStruct
	// These functions will be called when processing fields with default:"$(funcName)"
	conftagz.RegisterDefaultFunc("defaultFeatureFlags", defaultFeatureFlags)
	conftagz.RegisterDefaultFunc("defaultRuntimeInfo", defaultRuntimeInfo)

	// Get a struct with all defaults applied
	// The input struct is only used as a type template
	result, err := conftagz.GetDefaultOnlyStruct(&AppConfig{})
	if err != nil {
		log.Fatalf("Error getting defaults: %v", err)
	}

	// Type assert to our config type
	config := result.(*AppConfig)

	// Print the configuration with all defaults applied
	fmt.Println("Application Configuration (with defaults):")
	fmt.Println("==========================================")
	fmt.Printf("AppName: %s\n", config.AppName)
	fmt.Printf("Version: %s\n", config.Version)
	fmt.Printf("Debug: %v\n", config.Debug)
	fmt.Println()

	fmt.Println("Database Configuration:")
	fmt.Printf("  Host: %s\n", config.Database.Host)
	fmt.Printf("  Port: %d\n", config.Database.Port)
	fmt.Printf("  Username: %s\n", config.Database.Username)
	fmt.Printf("  Password: %s\n", config.Database.Password)
	fmt.Printf("  Database: %s\n", config.Database.Database)
	fmt.Printf("  SSLMode: %s\n", config.Database.SSLMode)
	fmt.Println()

	fmt.Println("Server Configuration:")
	fmt.Printf("  ListenAddr: %s\n", config.Server.ListenAddr)
	fmt.Printf("  Port: %d\n", config.Server.Port)
	fmt.Printf("  ReadTimeout: %d\n", config.Server.ReadTimeout)
	fmt.Printf("  WriteTimeout: %d\n", config.Server.WriteTimeout)
	fmt.Printf("  MaxHeaderBytes: %d\n", config.Server.MaxHeaderBytes)
	fmt.Println()

	fmt.Println("Log Configuration:")
	fmt.Printf("  Level: %s\n", config.Log.Level)
	fmt.Printf("  Format: %s\n", config.Log.Format)
	fmt.Printf("  Output: %v\n", config.Log.Output)
	fmt.Println()

	fmt.Println("Feature Flags (from custom function):")
	fmt.Printf("  EnableMetrics: %v\n", config.Features.EnableMetrics)
	fmt.Printf("  EnableProfiling: %v\n", config.Features.EnableProfiling)
	fmt.Printf("  EnableTracing: %v\n", config.Features.EnableTracing)
	fmt.Printf("  MaintenanceMode: %v\n", config.Features.MaintenanceMode)
	fmt.Printf("  MaxConcurrency: %d\n", config.Features.MaxConcurrency)
	fmt.Printf("  BuildInfo: %s\n", config.Features.BuildInfo)
	fmt.Println()

	fmt.Println("Runtime Info (from custom function with dynamic values):")
	fmt.Printf("  Hostname: %s\n", config.Runtime.Hostname)
	fmt.Printf("  StartTime: %s\n", config.Runtime.StartTime)
	fmt.Printf("  PID: %d\n", config.Runtime.PID)
}

func main() {
	RunMain()
}
