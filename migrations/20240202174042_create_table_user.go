package migrations

import (
	"gin-boilerplate/app/entity"

	"gorm.io/gorm"
)

type CREATE_TABLE_USER_MIGRATION struct{}

func (m *CREATE_TABLE_USER_MIGRATION) Migrate(db *gorm.DB) error {
	db.Migrator().CreateTable(&entity.User{})
	return nil
}

func (m *CREATE_TABLE_USER_MIGRATION) Rollback(db *gorm.DB) error {
	db.Migrator().DropTable(&entity.User{})
	return nil
}
