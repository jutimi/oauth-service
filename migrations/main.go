package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// Migration interface defines methods for applying and rolling back migrations.
type migration interface {
	Migrate(db *gorm.DB) error
	Rollback(db *gorm.DB) error
}

// MigrationList holds all registered migrations.
var MigrationList []migration

// RegisterMigration registers a migration to the list.
func RegisterMigration(migration migration) {
	MigrationList = append(MigrationList, migration)
}

// RegisterMigrations registers multiple migrations to the list.
func RegisterMigrations(migrations ...migration) {
	MigrationList = append(MigrationList, migrations...)
}

// GetMigrations returns all registered migrations.
func GetMigrations() []migration {
	return MigrationList
}

// RunMigrations applies all registered migrations to the database.
func RunMigrations(db *gorm.DB) error {
	for _, migration := range MigrationList {
		err := migration.Migrate(db)
		if err != nil {
			// Rollback all applied migrations
			rollbackErr := RollbackMigrations(db)
			if rollbackErr != nil {
				return fmt.Errorf("failed to rollback migrations after error: %w", rollbackErr)
			}

			return fmt.Errorf("failed to apply migration: %w", err)
		}
	}
	return nil
}

// RollbackMigrations rolls back all registered migrations from the database.
func RollbackMigrations(db *gorm.DB) error {
	for i := len(MigrationList) - 1; i >= 0; i-- {
		err := MigrationList[i].Rollback(db)
		if err != nil {
			return fmt.Errorf("failed to rollback migration: %w", err)
		}
	}
	return nil
}
