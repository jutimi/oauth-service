package migrations

import (
	"gorm.io/gorm"
)

type PREREQUISITES_MIGRATION struct{}

func (m *PREREQUISITES_MIGRATION) Migrate(db *gorm.DB) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	return nil
}

func (m *PREREQUISITES_MIGRATION) Rollback(db *gorm.DB) error {
	return nil
}
