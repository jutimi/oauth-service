package postgres_repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func findById(id uuid.UUID, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", field), id)
	}
}

func paginate(limit, offset int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit).Offset(offset)
	}
}

func findByText(text, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		searchText := "%" + text + "%"
		return db.Where(fmt.Sprintf("%s LIKE ?", field), searchText)
	}
}

func findBySlice[T []uuid.UUID | []string | []int](data T, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IN ?", field), data)
	}
}
