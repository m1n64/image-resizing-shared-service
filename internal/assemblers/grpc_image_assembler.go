package assemblers

import (
	images "image-resizing-shared/internal/delivery/grpc/pb"
	"image-resizing-shared/internal/domain"
	"image-resizing-shared/internal/ports"
	"image-resizing-shared/pkg/utils"
)

func GRPCBuildImage(image *ports.UploadResult) *images.ImageResponse {
	return &images.ImageResponse{
		Id:          image.ID,
		OriginalUrl: utils.BuildFileURL(image.OriginalKey),
		Status:      image.Status,
		Mime:        image.MIME,
		Size:        image.Size,
	}
}

func GRPCBuildImageWithThumbnails(image *domain.Image) *images.ImageResponse {
	var thumbnails []*images.ThumbnailShort
	for _, thumb := range image.Thumbnails {
		url := utils.BuildFileURL(thumb.Path)

		thumbnails = append(thumbnails, &images.ThumbnailShort{
			Size: thumb.Size,
			Url:  url,
			Type: thumb.Type,
		})
	}

	webpUrl := utils.BuildFileURL(image.WebPPath)

	return &images.ImageResponse{
		Id:            image.ID.String(),
		OriginalUrl:   utils.BuildFileURL(image.OriginalPath),
		CompressedUrl: &webpUrl,
		Status:        string(image.Status),
		Size:          image.Size,
		Mime:          image.MIME,
		ErrorMessage:  image.ErrorMessage,
		Thumbnails:    thumbnails,
	}
}
