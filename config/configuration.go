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

type JWT struct {
	Issuer              string `mapstructure:"issuer"`
	UserAccessTokenKey  string `mapstructure:"user_access_token_key"`
	UserRefreshTokenKey string `mapstructure:"user_refresh_token_key"`
}

type Server struct {
	Port int `mapstructure:"port"`
}

type Configuration struct {
	Database Database `mapstructure:"database"`
	Server   Server   `mapstructure:"server"`
	Jwt      JWT      `mapstructure:"jwt"`
}

func Init(filePath string) {
	var configuration *Configuration

	viper.SetConfigFile(filePath)
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	if err := viper.Unmarshal(&configuration); err != nil {
		log.Fatalf("error_decode_config: %v", err)
	}

	if configuration.Server.Port == 0 {
		configuration.Server.Port = 8080
	}

	config = configuration
}

func GetConfiguration() *Configuration {
	return config
}
