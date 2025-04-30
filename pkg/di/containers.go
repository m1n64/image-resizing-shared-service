package di

import (
	"gorm.io/gorm"
	"image-resizing-shared/internal/app"
	db2 "image-resizing-shared/internal/infrastructure/db"
	"image-resizing-shared/internal/ports"
	"image-resizing-shared/pkg/utils"
	"log"
	"os"
	"path/filepath"
)

const (
	DbFile     = "data.db"
	UploadsDir = "uploads"
)

type Dependencies struct {
	DB         *gorm.DB
	UploadsDir string
	// Repositories
	ImageRepository     ports.ImageRepository
	ThumbnailRepository ports.ThumbnailRepository
	// Usecases
	ResizeService ports.ResizeUseCase
	ImageService  ports.ImageUseCase
}

func InitDependencies() *Dependencies {
	_, uploadsDir, dbFile := initFS()

	db := utils.InitDB(dbFile)
	utils.AutoMigrate(db)

	// Repositories
	imageRepository := db2.NewImageRepository(db)
	thumbnailRepository := db2.NewThumbnailRepository(db)

	// UseCases
	resizeService := app.NewResizeService(db, thumbnailRepository, imageRepository, uploadsDir)
	imageService := app.NewImageService(db, imageRepository, thumbnailRepository, resizeService, uploadsDir)

	return &Dependencies{
		DB:                  db,
		UploadsDir:          uploadsDir,
		ImageRepository:     imageRepository,
		ThumbnailRepository: thumbnailRepository,
		ResizeService:       resizeService,
		ImageService:        imageService,
	}
}

func initFS() (string, string, string) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = utils.GetDataPath(DbFile)
	}

	dataDir := filepath.Dir(dbPath)

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("error creating data dirrectory: %v", err)
	}

	uploadsDir := filepath.Join(dataDir, UploadsDir)

	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Fatalf("error creating data dirrectory: %v", err)
	}

	return dataDir, uploadsDir, dbPath
}
