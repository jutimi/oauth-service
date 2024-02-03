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
	db.Exec("CREATE TABLE migrations (id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), name TEXT NOT NULL, file TEXT NOT NULL);")
	return nil
}

func downPrerequisites(db *gorm.DB) error {
	return nil
}
