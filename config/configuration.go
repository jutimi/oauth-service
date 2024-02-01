package config

import (
	"log"

	"github.com/spf13/viper"
)

var config *Configuration

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type Server struct {
	Port int `mapstructure:"port"`
}

type Configuration struct {
	Database Database `mapstructure:"database"`
	Server   Server   `mapstructure:"server"`
}

func Init(filePath string) {
	var configuration *Configuration

	viper.SetConfigFile(filePath)
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	if err := viper.Unmarshal(&configuration); err != nil {
		log.Fatalf("error_decode_config: %v", err)
	}

	config = configuration
}

func GetConfiguration() *Configuration {
	return config
}
