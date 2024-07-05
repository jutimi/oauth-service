package migrations

import (
	"oauth-server/app/entity"

	"gorm.io/gorm"
)

func init() {
	RegisterUpFunc("upCreateTableUser", upCreateTableUser)
	RegisterDownFunc("downCreateTableUser", downCreateTableUser)
}

func upCreateTableUser(db *gorm.DB) error {
	if !db.Migrator().HasTable(&entity.User{}) {
		db.Migrator().CreateTable(&entity.User{})
	}

	return nil
}

func downCreateTableUser(db *gorm.DB) error {
	db.Migrator().DropTable(&entity.User{})
	return nil
}
