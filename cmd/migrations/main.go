package main

import (
	"fmt"
	"gin-boilerplate/config"
	"gin-boilerplate/migrations"
	"gin-boilerplate/package/database"
	"gin-boilerplate/utils"
	"os"
	"time"
)

const (
	ACTION_CREATE = "create"
	ACTION_UP     = "up"
	ACTION_DOWN   = "down"
)

func init() {
	rootDir := utils.RootDir()
	configFile := fmt.Sprintf("%s/config.yml", rootDir)

	config.Init(configFile)
	database.Init()
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Invalid arguments")
		return
	}

	action := args[1]
	switch action {
	case ACTION_CREATE:
		createMigration(args)
	case ACTION_UP:
		upMigration(args)
	case ACTION_DOWN:
		downMigration(args)
	default:
		fmt.Println("Action not supported")
	}
}

func createMigration(args []string) {
	if len(args) < 2 {
		fmt.Println("Missing file name")
		return
	}

	rootDir := utils.RootDir()

	name := args[2]
	currentTime := time.Now().Format(utils.TIME_STAMP_FORMAT)
	fileName := fmt.Sprintf("%s_%s.go", currentTime, name)
	filePath := fmt.Sprintf("%s/migrations/%s", rootDir, fileName)
	structName := utils.ConvertToUppercase(name + "_Migration")

	fileContent := fmt.Sprintf(`
package migrations

import (
	"gorm.io/gorm"
)

type %s struct{}

func (m *%s) Migrate(db *gorm.DB) error {
	return nil
}

func (m *%s) Rollback(db *gorm.DB) error {
	return nil
}
	`, structName, structName, structName)

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating migration file:", err.Error())
		return
	}
	defer file.Close()

	_, err = file.WriteString(fileContent)
	if err != nil {
		fmt.Println("Error generate migration file:", err.Error())
	}
}

func upMigration(args []string) {
	db := database.GetDB()

	// Manually register the migration
	prerequisitesMigration := &migrations.PREREQUISITES_MIGRATION{}
	createTableMigration := &migrations.CREATE_TABLE_USER_MIGRATION{}
	migrations.RegisterMigrations(
		prerequisitesMigration,
		createTableMigration,
	)

	// Run all registered migrations
	if err := migrations.RunMigrations(db); err != nil {
		fmt.Println("Failed to run migrations:", err)
		return
	}
}

func downMigration(args []string) {
	return
}
