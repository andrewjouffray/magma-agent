package config

import (
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {
	// Create a temporary config file
	tempFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data to the temp file
	configData := `
device_id: test-device-id
`
	if _, err := tempFile.Write([]byte(configData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Read the config using the ReadConfig function
	config, err := ReadConfig(tempFile.Name())
	if err != nil {
		t.Fatalf("ReadConfig returned an error: %v", err)
	}

	// Verify the config values
	expectedDeviceID := "test-device-id"
	if config.DeviceID != expectedDeviceID {
		t.Errorf("Expected DeviceID %s, got %s", expectedDeviceID, config.DeviceID)
	}
}

func TestReadConfig_FileNotFound(t *testing.T) {
	_, err := ReadConfig("nonexistent.yaml")
	if err == nil {
		t.Fatal("Expected an error for nonexistent file, but got nil")
	}
}

func TestReadConfig_InvalidYAML(t *testing.T) {
	// Create a temporary config file with invalid YAML
	tempFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write invalid YAML data to the temp file
	invalidConfigData := `
device_id: test-device-id
invalid_yaml
`
	if _, err := tempFile.Write([]byte(invalidConfigData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Attempt to read the config using the ReadConfig function
	_, err = ReadConfig(tempFile.Name())
	if err == nil {
		t.Fatal("Expected an error for invalid YAML, but got nil")
	}
}
