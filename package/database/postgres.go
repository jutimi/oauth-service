package database

import (
	"fmt"
	"gin-boilerplate/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
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

	db = conn
}

func GetDB() *gorm.DB {
	return db
}
