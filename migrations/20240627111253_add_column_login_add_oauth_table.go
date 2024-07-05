package migrations

import (
	"oauth-server/app/entity"

	"gorm.io/gorm"
)

func init() {
	RegisterUpFunc("upAddColumnLoginAddOauthTable", upAddColumnLoginAddOauthTable)
	RegisterDownFunc("downAddColumnLoginAddOauthTable", downAddColumnLoginAddOauthTable)
}

func upAddColumnLoginAddOauthTable(db *gorm.DB) error {
	if db.Migrator().HasTable(&entity.Oauth{}) {
		if err := db.Migrator().AddColumn(&entity.Oauth{}, "login_at"); err != nil {
			return err
		}
	}

	return nil
}

func downAddColumnLoginAddOauthTable(db *gorm.DB) error {
	if db.Migrator().HasTable(&entity.Oauth{}) {
		if err := db.Migrator().DropColumn(&entity.Oauth{}, "login_at"); err != nil {
			return err
		}
	}

	return nil
}
