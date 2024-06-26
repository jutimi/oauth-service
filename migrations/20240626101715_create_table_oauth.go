package migrations

import (
	"gin-boilerplate/app/entity"

	"gorm.io/gorm"
)

func init() {
	RegisterUpFunc("upCreateTableOauth", upCreateTableOauth)
	RegisterDownFunc("downCreateTableOauth", downCreateTableOauth)
}

func upCreateTableOauth(db *gorm.DB) error {
	if !db.Migrator().HasTable(&entity.Oauth{}) {
		db.Migrator().CreateTable(&entity.Oauth{})
	}

	return nil
}

func downCreateTableOauth(db *gorm.DB) error {
	db.Migrator().DropTable(&entity.Oauth{})
	return nil
}
