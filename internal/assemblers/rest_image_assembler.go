package assemblers

import (
	"image-resizing-shared/internal/domain"
	"image-resizing-shared/internal/dto"
	"image-resizing-shared/internal/ports"
	"image-resizing-shared/pkg/utils"
)

func BuildImage(image *ports.UploadResult) *dto.ImageWithThumbnails {
	return &dto.ImageWithThumbnails{
		ID:          image.ID,
		OriginalUrl: utils.BuildFileURL(image.OriginalKey),
		Status:      image.Status,
		Size:        image.Size,
		MIME:        image.MIME,
	}
}

func BuildImageWithThumbnails(image *domain.Image) *dto.ImageWithThumbnails {
	thumbnails := make([]dto.ThumbnailShort, 0, len(image.Thumbnails))
	for _, thumb := range image.Thumbnails {
		thumbnails = append(thumbnails, dto.ThumbnailShort{
			Size: thumb.Size,
			Url:  utils.BuildFileURL(thumb.Path),
			Type: thumb.Type,
		})
	}

	return &dto.ImageWithThumbnails{
		ID:            image.ID.String(),
		OriginalUrl:   utils.BuildFileURL(image.OriginalPath),
		CompressedUrl: utils.BuildFileURL(image.WebPPath),
		Status:        string(image.Status),
		Size:          image.Size,
		MIME:          image.MIME,
		ErrorMessage:  image.ErrorMessage,
		Thumbnails:    thumbnails,
	}
}
