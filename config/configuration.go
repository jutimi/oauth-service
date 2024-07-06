package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

var config *Configuration

type Configuration struct {
	MysqlDB    MysqlDatabase    `mapstructure:"mysql"`
	PostgresDB PostgresDatabase `mapstructure:"postgres"`
	Server     Server           `mapstructure:"server"`
	Jwt        JWT              `mapstructure:"jwt"`
	GRPC       GRPC             `mapstructure:"grpc"`
}

func Init(filePath string) {
	var configuration *Configuration

	// Check file exist base on filepath
	if _, err := os.Stat(filePath); err != nil {
		log.Fatalf("File Config Not Found: %s", err.Error())
	}

	viper.SetConfigFile(filePath)
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	if err := viper.Unmarshal(&configuration); err != nil {
		log.Fatalf("Error Decode Config: %s", err.Error())
	}

	if configuration.Server.Port == 0 {
		configuration.Server.Port = 8080
	}

	config = configuration
	log.Println("Init Config Success!")
}

func GetConfiguration() *Configuration {
	return config
}
