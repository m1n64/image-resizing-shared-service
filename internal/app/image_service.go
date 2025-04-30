package app

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"image-resizing-shared/internal/domain"
	"image-resizing-shared/internal/ports"
	"image-resizing-shared/pkg/utils"
	"log"
	"os"
	"path/filepath"
)

type ImageService struct {
	db                  *gorm.DB
	imageRepository     ports.ImageRepository
	thumbnailRepository ports.ThumbnailRepository
	resizeService       ports.ResizeUseCase
	storageRoot         string
}

func NewImageService(
	db *gorm.DB,
	imageRepo ports.ImageRepository,
	thumbRepo ports.ThumbnailRepository,
	resizeService ports.ResizeUseCase,
	storageRoot string,
) ports.ImageUseCase {
	return &ImageService{
		db:                  db,
		imageRepository:     imageRepo,
		thumbnailRepository: thumbRepo,
		resizeService:       resizeService,
		storageRoot:         storageRoot,
	}
}

func (s *ImageService) UploadOriginal(ctx context.Context, filePath string, contentType string) (*ports.UploadResult, error) {
	id := uuid.New()

	ext := filepath.Ext(filePath)
	fileName := filepath.Join("originals", id.String()+ext)
	originalPath := filepath.Join(s.storageRoot, fileName)

	if err := utils.SaveFileToStorage(filePath, originalPath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	fileSize := utils.FileSize(filePath)

	image := &domain.Image{
		ID:           id,
		Filename:     filepath.Base(filePath),
		MIME:         contentType,
		Size:         fileSize,
		OriginalPath: fileName,
		Status:       domain.StatusPending,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		imgRepo := s.imageRepository.WithTx(tx)
		return imgRepo.Save(image)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save image record: %w", err)
	}

	go s.compressAndDispatch(ctx, id, originalPath)

	return &ports.UploadResult{
		ID:          id.String(),
		OriginalKey: fileName,
		Status:      string(domain.StatusPending),
		MIME:        contentType,
		Size:        fileSize,
	}, nil
}

func (s *ImageService) FindByID(ctx context.Context, id string) (*domain.Image, error) {
	return s.imageRepository.FindByID(id)
}

func (s *ImageService) compressAndDispatch(ctx context.Context, id uuid.UUID, originalKey string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic: %v", r)
		}
	}()

	originalFile, err := utils.ReadFileFromStorage(originalKey)
	if err != nil {
		s.markAsError(ctx, id, fmt.Errorf("failed to get file as bytes: %w", err))
		return
	}

	webpPath, err := utils.ConvertBytesToWebp(originalFile)
	if err != nil {
		s.markAsError(ctx, id, fmt.Errorf("failed to convert to webp: %w", err))
		return
	}
	defer os.Remove(webpPath)

	webpName := filepath.Join("webp", id.String()+".webp")
	webpFilePath := filepath.Join(s.storageRoot, webpName)

	if err := utils.SaveFileToStorage(webpPath, webpFilePath); err != nil {
		s.markAsError(ctx, id, fmt.Errorf("failed to upload compressed webp: %w", err))
		return
	}

	image, err := s.imageRepository.FindByID(id.String())
	if err != nil {
		log.Printf("failed to find image %s: %v", id, err)
		return
	}

	image.WebPPath = webpName
	image.Status = domain.StatusProcessing

	if err := s.imageRepository.Update(image); err != nil {
		log.Printf("failed to update image %s: %v", image.ID, err)
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recovered from panic: %v", r)
			}
		}()

		if err := s.resizeService.ResizeThumbnails(context.Background(), id, webpFilePath); err != nil {
			log.Printf("failed to resize thumbnails: %v", err)
			s.markAsError(context.Background(), id, err)
		}
	}()
}

func (s *ImageService) markAsError(ctx context.Context, id uuid.UUID, originalErr error) {
	image, err := s.imageRepository.FindByID(id.String())
	if err != nil {
		return
	}

	errorMessage := originalErr.Error()
	image.Status = domain.StatusError
	image.ErrorMessage = &errorMessage

	if err := s.imageRepository.Update(image); err != nil {
		log.Printf("failed to update image %s: %v", image.ID, err)
	}
}
