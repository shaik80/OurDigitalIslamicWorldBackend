package config

import (
	"sync"

	"github.com/spf13/viper"
)

// Config holds the configuration values for the application
type Config struct {
	App        AppConfig     `yaml:"app"`
	OpenSearch OpenSearch    `yaml:"opensearch"`
	Server     ServerConfig  `yaml:"server"`
	Logging    LoggingConfig `yaml:"logging"`
}

// AppConfig holds information about the application
type AppConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// OpenSearch holds the configuration for connecting to OpenSearch
type OpenSearch struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// ServerConfig holds the configuration values for the web server
type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// LoggingConfig holds the configuration for logging
type LoggingConfig struct {
	LogLevel string `yaml:"log_level"`
}

var (
	appConfig     Config
	appConfigOnce sync.Once
)

// Load initializes the configuration once
func Load() error {
	var err error
	appConfigOnce.Do(func() {
		err = loadConfig()
	})
	return err
}

// loadConfig loads the configuration from the file
func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set default values
	viper.SetDefault("app.name", "MyFiberApp")
	viper.SetDefault("app.description", "A sample Fiber application")

	viper.SetDefault("opensearch.host", "localhost")
	viper.SetDefault("opensearch.port", "9200")
	viper.SetDefault("opensearch.username", "admin")
	viper.SetDefault("opensearch.password", "admin")

	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", "8080")

	viper.SetDefault("logging.log_level", "info")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&appConfig); err != nil {
		return err
	}

	return nil
}

// GetConfig returns the loaded configuration
func GetConfig() Config {
	return appConfig
}

// GetString gets a string configuration value by key
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt gets an integer configuration value by key
func GetInt(key string) int {
	return viper.GetInt(key)
}
