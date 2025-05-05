package config

import (
	"log"
)

// Config holds the server configuration
type Config struct {
	Port       string
	StorageDir string
}

// GlobalConfig holds the globally accessible config
var GlobalConfig *Config

// LoadConfig initializes the config from command-line flags
func LoadConfig(port, storageDir string) *Config {
	// Assign the value to the global config variable
	GlobalConfig = &Config{
		Port:       port,
		StorageDir: storageDir,
	}

	// Log the configuration for debugging purposes
	log.Printf("Server will start on port %s", GlobalConfig.Port)

	return GlobalConfig
}

// GetStorage returns the storage directory from the global config
func GetStorage() string {
	return GlobalConfig.StorageDir
}
