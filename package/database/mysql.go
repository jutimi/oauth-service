package database

import (
	"fmt"
	"log"
	"oauth-server/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mysqlDB *gorm.DB

func InitMysql() {
	config := config.GetConfiguration().MysqlDB

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error_connecting_to_database: %v", err)
	}

	postgresDB = conn
}

func GetMysql() *gorm.DB {
	return mysqlDB
}

func BeginMysqlTransaction() *gorm.DB {
	return GetMysql().Begin()
}
