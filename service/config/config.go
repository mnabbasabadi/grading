// Package config provides a configuration struct and a function to read the config file
package config

import (
	"fmt"

	"github.com/mnabbasabadi/grading/service/foundation/db"
	"golang.org/x/exp/slog"

	"github.com/spf13/viper"
)

// Config is a struct that holds the configuration values
type (
	// DB ...
	DB struct {
		User     string
		Password string
		Host     string
		Port     string
		DBName   string
		SslMode  db.SSLMode
	}

	// Config is a struct that holds the configuration values
	Config struct {
		Host        string
		Port        string
		LogLevel    string
		ServiceName string
		DB          DB
	}
)

func initConfig() error {

	// Bind environment variables to config keys
	viper.AutomaticEnv()

	configFile := viper.GetString("CONFIG_FILE")
	// Read the config file (if it exists)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		// Handle the error if the config file doesn't exist (you can choose to create it or use only default values)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %v", err)
		}
	}

	keys := []string{
		"Host", "Port", "LogLevel", "ServiceName",
		"DB.User", "DB.Password", "DB.Host", "DB.Port", "DB.DBName", "DB.Sslmode",
	}
	if err := bindEnv(keys...); err != nil {
		return fmt.Errorf("failed to bind environment variables: %v", err)
	}

	return nil
}

func bindEnv(keys ...string) error {
	for _, key := range keys {
		if err := viper.BindEnv(key); err != nil {
			return fmt.Errorf("failed to bind environment variable %s: %v", key, err)
		}
	}
	return nil
}

// NewConfig reads the config file and returns a Config struct
func NewConfig() (*Config, error) {
	if err := initConfig(); err != nil {
		return nil, fmt.Errorf("failed to initialize config: %v", err)
	}
	// Unmarshal the config into a struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}

// RecoverAndLogPanic ...
func RecoverAndLogPanic(logger *slog.Logger) {
	if err := recover(); err != nil {
		logger.With(err).Error("panic")
	}
}
