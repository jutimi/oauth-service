package main

import (
	"fmt"
	"oauth-server/config"
	"oauth-server/migrations"
	"oauth-server/package/database"
	"oauth-server/utils"
	"os"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

func init() {
	rootDir := utils.RootDir()
	configFile := fmt.Sprintf("%s/config.yml", rootDir)

	config.Init(configFile)
	database.InitPostgres()
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
	db := database.GetPostgres()

	if len(args) > 3 {
		fmt.Println("Too many arguments")
		return
	}

	fileName := args[2]
	if fileName == "" {
		runMultipleMigrations(action, db)
		return
	}

	runMigration(fileName, action, db)
}

// -------------------------------------------------------------------------------
func generateFuncName(fileName, action string) string {
	name := generateName(fileName)
	switch action {
	case migrations.ACTION_UP, migrations.ACTION_DOWN:
		return action + utils.ConvertToCamelCase(utils.RemoveFileNameExtension(name))
	default:
		return ""
	}
}

func generateName(fileName string) string {
	parts := strings.Split(fileName, "_")
	nameParts := parts[1:]
	name := strings.Join(nameParts, "_")
	return name
}

func logMigration(name, fileName string) error {
	migration := migrations.MigrationTable{
		Name: name,
		File: fileName,
	}
	result := database.GetPostgres().Create(&migration)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func getMigrationsByName(name string) ([]migrations.MigrationTable, error) {
	var data []migrations.MigrationTable

	db := database.GetPostgres()
	result := db.Where(&migrations.MigrationTable{Name: name}).Find(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

func removeMigrationByName(name string) error {
	var data *migrations.MigrationTable

	db := database.GetPostgres()
	result := db.Where(&migrations.MigrationTable{Name: name}).First(&data)
	if result.Error != nil {
		return result.Error
	}

	result = db.Delete(&data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func runMultipleMigrations(action string, db *gorm.DB) {
	// Get all migrations file
	dirPath := utils.RootDir() + "/migrations"
	files, err := utils.GetAllFileInDir(dirPath)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Define the regular expression pattern
		pattern := `^\d{14}_\w+\.go$`

		// Compile the regular expression
		regex, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Println("Error compiling regular expression:", err)
			return
		}
		match := regex.MatchString(file.Name())
		if !match || strings.TrimSpace(file.Name()) == " " {
			fmt.Println("File name does not match the pattern:", file.Name())
			continue
		}

		fmt.Println("Running migration:", file.Name())
		runMigration(file.Name(), action, db)
		fmt.Println("Finish migration:", file.Name())
	}
}

func runMigration(fileName, action string, db *gorm.DB) {
	name := generateName(fileName)

	// Validate migration
	if action == migrations.ACTION_UP && !strings.Contains(fileName, "prerequisites") {
		migrationsData, err := getMigrationsByName(name)
		if err != nil {
			fmt.Printf("Error getting migrations %s by name: %s \n", fileName, err.Error())
			return
		}
		if len(migrationsData) > 0 {
			fmt.Printf("Migration %s already run \n", fileName)
			return
		}
	}

	// Run migration file
	funcName := generateFuncName(fileName, action)
	if err := migrations.Run(funcName, action, db); err != nil {
		fmt.Printf("Error running migrations %s: %s \n", fileName, err.Error())
		return
	}

	// Handle logging in database
	switch action {
	case migrations.ACTION_UP:
		// Log migration
		fmt.Println("Running migration:", fileName)
		if err := logMigration(name, fileName); err != nil {
			fmt.Printf("Error logging %s migration: %s \n", fileName, err.Error())
			return
		}
		fmt.Println("Finish migration:", fileName)
	case migrations.ACTION_DOWN:
		removeMigrationByName(name)
	default:
		break
	}
}
