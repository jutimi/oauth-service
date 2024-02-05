package migrations

import (
	"gorm.io/gorm"
)

func init() {
	RegisterUpFunc("upPrerequisites", upPrerequisites)
	RegisterDownFunc("downPrerequisites", downPrerequisites)
}

func upPrerequisites(db *gorm.DB) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if !db.Migrator().HasTable(&MigrationTable{}) {
		db.Migrator().CreateTable(&MigrationTable{})
	}
	return nil
}

func downPrerequisites(db *gorm.DB) error {
	return nil
}
