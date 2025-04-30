package di

import (
	"gorm.io/gorm"
	"image-resizing-shared/internal/app"
	db2 "image-resizing-shared/internal/infrastructure/db"
	"image-resizing-shared/internal/ports"
	"image-resizing-shared/pkg/config"
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
	Config     *config.Config
	// Repositories
	ImageRepository     ports.ImageRepository
	ThumbnailRepository ports.ThumbnailRepository
	// Usecases
	ResizeService ports.ResizeUseCase
	ImageService  ports.ImageUseCase
}

func InitDependencies(config *config.Config) *Dependencies {
	_, uploadsDir, dbFile := initFS(config)

	db := utils.InitDB(dbFile)
	utils.AutoMigrate(db)

	// Repositories
	imageRepository := db2.NewImageRepository(db)
	thumbnailRepository := db2.NewThumbnailRepository(db)

	// UseCases
	resizeService := app.NewResizeService(db, thumbnailRepository, imageRepository, uploadsDir, &config.Compression)
	imageService := app.NewImageService(db, imageRepository, thumbnailRepository, resizeService, uploadsDir, &config.Compression)

	return &Dependencies{
		DB:                  db,
		UploadsDir:          uploadsDir,
		Config:              config,
		ImageRepository:     imageRepository,
		ThumbnailRepository: thumbnailRepository,
		ResizeService:       resizeService,
		ImageService:        imageService,
	}
}

func initFS(config *config.Config) (string, string, string) {
	dbPath := config.DBPath
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
