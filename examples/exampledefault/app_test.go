package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.izuma.io/conftagz"
)

func TestExampleDefaultExecution(t *testing.T) {
	// Register custom default functions (same as RunMain does)
	conftagz.RegisterDefaultFunc("defaultFeatureFlags", defaultFeatureFlags)
	conftagz.RegisterDefaultFunc("defaultRuntimeInfo", defaultRuntimeInfo)

	// Test that GetDefaultOnlyStruct works correctly with our AppConfig
	result, err := conftagz.GetDefaultOnlyStruct(&AppConfig{})
	assert.NoError(t, err)
	assert.NotNil(t, result)

	config := result.(*AppConfig)

	// Verify top-level defaults
	assert.Equal(t, "MyApplication", config.AppName)
	assert.Equal(t, "1.0.0", config.Version)
	assert.Equal(t, false, config.Debug)

	// Verify nested pointer struct (Database) was created and has defaults
	assert.NotNil(t, config.Database)
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, 5432, config.Database.Port)
	assert.Equal(t, "postgres", config.Database.Username)
	assert.Equal(t, "secret", config.Database.Password)
	assert.Equal(t, "myapp", config.Database.Database)
	assert.Equal(t, "disable", config.Database.SSLMode)

	// Verify embedded struct (Server) has defaults
	assert.Equal(t, "0.0.0.0", config.Server.ListenAddr)
	assert.Equal(t, 8080, config.Server.Port)
	assert.Equal(t, 30, config.Server.ReadTimeout)
	assert.Equal(t, 30, config.Server.WriteTimeout)
	assert.Equal(t, 1048576, config.Server.MaxHeaderBytes)

	// Verify nested pointer struct (Log) was created and has defaults
	assert.NotNil(t, config.Log)
	assert.Equal(t, "info", config.Log.Level)
	assert.Equal(t, "json", config.Log.Format)
	assert.Equal(t, []string{"stdout", "file"}, config.Log.Output)

	// Verify custom function defaults (Features)
	assert.NotNil(t, config.Features)
	assert.Equal(t, true, config.Features.EnableMetrics)
	assert.Equal(t, false, config.Features.EnableProfiling)
	assert.Equal(t, true, config.Features.EnableTracing)
	assert.Equal(t, false, config.Features.MaintenanceMode)
	assert.Equal(t, 100, config.Features.MaxConcurrency)
	assert.Equal(t, "v1.0.0-abc123", config.Features.BuildInfo)

	// Verify custom function defaults (Runtime) - dynamic values
	assert.NotNil(t, config.Runtime)
	assert.NotEmpty(t, config.Runtime.Hostname)
	assert.NotEmpty(t, config.Runtime.StartTime)
	assert.NotZero(t, config.Runtime.PID)
}

func TestRunMain(t *testing.T) {
	// Just ensure RunMain doesn't panic
	assert.NotPanics(t, func() {
		RunMain()
	})
}
