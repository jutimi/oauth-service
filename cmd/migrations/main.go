package main

import (
	"fmt"
	"gin-boilerplate/config"
	"gin-boilerplate/migrations"
	"gin-boilerplate/package/database"
	"gin-boilerplate/utils"
	"os"
	"strings"
	"time"
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
	case migrations.ACTION_CREATE:
		createMigration(args)
	case migrations.ACTION_UP, migrations.ACTION_DOWN:
		migrate(args, action)
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
	upFunc := "up" + utils.ConvertToCamelCase(name)
	downFunc := "down" + utils.ConvertToCamelCase(name)

	fileContent := fmt.Sprintf(`
package migrations

import (
	"gorm.io/gorm"
)

func init() {
	RegisterUpFunc("%s", %s)
	RegisterDownFunc("%s", %s)
}

func %s(db *gorm.DB) error {
	return nil
}

func %s(db *gorm.DB) error {
	return nil
}
	`,
		upFunc, upFunc,
		downFunc, downFunc,
		upFunc,
		downFunc,
	)

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

func migrate(args []string, action string) {
	db := database.GetDB()

	if len(args) > 3 {
		fmt.Println("Too many arguments")
		return
	}

	fileName := args[2]
	if fileName == "" {
		return
	}

	funcName := generateFuncName(fileName, action)
	if err := migrations.Run(funcName, action, db); err != nil {
		fmt.Println("Error running migrations:", err.Error())
		return
	}

	return
}

func generateFuncName(fileName, action string) string {
	parts := strings.Split(fileName, "_")
	name := parts[1:]
	fileName = strings.Join(name, "_")
	switch action {
	case migrations.ACTION_UP, migrations.ACTION_DOWN:
		return action + utils.ConvertToCamelCase(utils.RemoveFileNameExtension(fileName))
	default:
		return ""
	}
}
