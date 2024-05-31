package database

import (
	"fmt"
	"gin-boilerplate/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresDB *gorm.DB

func InitPostgres() {
	config := config.GetConfiguration().Database

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		config.Host,
		config.User,
		config.Password,
		config.Database,
		config.Port,
	)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error_connecting_to_database: %v", err)
	}

	postgresDB = conn
}

func GetPostgres() *gorm.DB {
	return postgresDB
}

func BeginPostgresTransaction() *gorm.DB {
	return GetPostgres().Begin()
}
