package migrations

import "gorm.io/gorm"

type MigrationTable struct {
	gorm.Model
	ID   string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name string `gorm:"type:text;not null"`
	File string `gorm:"type:text;not null"`
}

func (MigrationTable) TableName() string {
	return "migrations"
}
