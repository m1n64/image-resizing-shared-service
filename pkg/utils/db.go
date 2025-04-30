package utils

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"image-resizing-shared/internal/domain"
	"path/filepath"
)

func InitDB(dbPath string) *gorm.DB {
	dsn := filepath.ToSlash(dbPath) + "?_journal_mode=WAL&_busy_timeout=5000"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Fail to open SQLite database: %v", err))
	}

	return db
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&domain.Image{}, &domain.Thumbnail{})
}
