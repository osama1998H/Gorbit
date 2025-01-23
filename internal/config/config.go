// internal/config/config.go
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port  int    `mapstructure:"port"`
		Host  string `mapstructure:"host"`
		Debug bool   `mapstructure:"debug"`
	} `mapstructure:"server"`

	App struct {
		Name    string `mapstructure:"name"`
		Version string `mapstructure:"version"`
		Env     string `yaml:"env"`
		APIKey  string `mapstructure:"api_key"`
		JWTSecret string `mapstructure:"jwt_secret"`
	} `mapstructure:"app"`

	Databases struct {
		MySQL struct {
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
			Database string `mapstructure:"database"`
		} `mapstructure:"mysql"`

		Postgres struct {
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
			Database string `mapstructure:"database"`
		} `mapstructure:"postgres"`

		MongoDB struct {
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
			Database string `mapstructure:"database"`
			AuthSource string `mapstructure:"auth_source"`
			AuthMechanism string `mapstructure:"auth_mechanism"`
		} `mapstructure:"mongodb"`
	} `mapstructure:"databases"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
}

func LoadConfig() (*Config, error) {
	// Set the configuration file name and path
	configPath := "configs"
	
	// Initialize Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	// Read the primary configuration
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Read database configuration
	databaseViper := viper.New()
	databaseViper.SetConfigName("database")
	databaseViper.SetConfigType("yaml")
	databaseViper.AddConfigPath(configPath)
	if err := databaseViper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading database config file: %w", err)
	}

	// Read Redis configuration
	redisViper := viper.New()
	redisViper.SetConfigName("redis")
	redisViper.SetConfigType("yaml")
	redisViper.AddConfigPath(configPath)
	if err := redisViper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading redis config file: %w", err)
	}

	// Merge configurations
	for k, v := range databaseViper.AllSettings() {
		viper.Set(k, v)
	}
	for k, v := range redisViper.AllSettings() {
		viper.Set(k, v)
	}

	// Unmarshal configuration
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &config, nil
}