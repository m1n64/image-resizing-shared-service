package handlers

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"image-resizing-shared/internal/assemblers"
	images "image-resizing-shared/internal/delivery/grpc/pb"
	"image-resizing-shared/internal/ports"
	"image-resizing-shared/pkg/utils"
)

type ImageGRPCHandler struct {
	images.UnimplementedImageServiceServer
	imageUseCase ports.ImageUseCase
}

func NewImageGRPCHandler(imageUseCase ports.ImageUseCase) *ImageGRPCHandler {
	return &ImageGRPCHandler{
		imageUseCase: imageUseCase,
	}
}

func (h *ImageGRPCHandler) UploadImage(ctx context.Context, req *images.UploadImageRequest) (*images.ImageResponse, error) {
	if req.Data == nil || len(req.Data) == 0 {
		return nil, status.Error(codes.InvalidArgument, "image data is required")
	}

	tempFilePath, contentType, err := utils.SaveTempFile(bytes.NewReader(req.Data))
	if err != nil {
		return nil, err
	}
	defer utils.RemoveFile(tempFilePath)

	result, err := h.imageUseCase.UploadOriginal(ctx, tempFilePath, contentType)
	if err != nil {
		return nil, err
	}

	return assemblers.GRPCBuildImage(result), nil
}

func (h *ImageGRPCHandler) GetImage(ctx context.Context, req *images.GetImageRequest) (*images.ImageResponse, error) {
	if uuid.Validate(req.Id) != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	imageData, err := h.imageUseCase.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return assemblers.GRPCBuildImageWithThumbnails(imageData), nil
}
