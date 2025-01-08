package config

import (
	"bufio"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// this is meant as a package for all the configuration constants
// variable configuration for the magma application will be handled by a separate readable config file
const (
	AppName      = "magma"
	Version      = "1.0.0"
	AppRoot      = "/etc/magma"
	TrackFile    = "/etc/magma/track"
	IgnoreFile   = "/etc/magma/ignore"
	SnapshotsDir = "/etc/magma/snapshots"
	ConfigFile   = "/etc/magma/config.yaml"
)

// VariableConfig holds the dynamically loaded configuration
var VariableConfig variableConfig

// variableConfig defines the structure of the YAML configuration
type variableConfig struct {
	DeviceID string `yaml:"device_id"`
}

// init initializes the package by reading the configuration file
func init() {
	var err error
	VariableConfig, err = ReadConfig(ConfigFile)
	if err != nil {
		// just notify that the system has not been initialized
		log.Println("Config file not found, please run 'magma init' if you aren't already")
	}
}

// ReadConfig reads the config.yaml file and unmarshals it into the variableConfig struct
func ReadConfig(configFile string) (variableConfig, error) {
	// Open the config file
	file, err := os.Open(configFile)
	if err != nil {
		return variableConfig{}, err
	}
	defer file.Close()

	// Decode the YAML file
	reader := bufio.NewReader(file)
	var config variableConfig
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&config); err != nil {
		return variableConfig{}, err
	}

	return config, nil
}
