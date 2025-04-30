package handlers

import (
	"github.com/gin-gonic/gin"
	"image-resizing-shared/internal/assemblers"
	"image-resizing-shared/internal/ports"
	"image-resizing-shared/pkg/utils"
	"net/http"
	"time"
)

var allowedTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

type ImageHandler struct {
	imageUseCase ports.ImageUseCase
}

func NewImageHandler(imageUseCase ports.ImageUseCase) *ImageHandler {
	return &ImageHandler{
		imageUseCase: imageUseCase,
	}
}

func (h *ImageHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "pong",
		"timestamp": time.Now(),
	})
}

func (h *ImageHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported content type"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer src.Close()

	tempFilePath, _, err := utils.SaveTempFile(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save temp file"})
		return
	}
	defer utils.RemoveFile(tempFilePath)

	result, err := h.imageUseCase.UploadOriginal(c.Request.Context(), tempFilePath, file.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assemblers.BuildImage(result))
}

func (h *ImageHandler) UploadImageBinary(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")
	if contentType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing Content-Type header"})
		return
	}

	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported content type"})
		return
	}

	tempFilePath, _, err := utils.SaveTempFile(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save temp file"})
		return
	}
	defer utils.RemoveFile(tempFilePath)

	result, err := h.imageUseCase.UploadOriginal(c.Request.Context(), tempFilePath, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assemblers.BuildImage(result))
}

func (h *ImageHandler) GetImage(c *gin.Context) {
	id := c.Param("id")

	image, err := h.imageUseCase.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if image == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}
	
	c.JSON(http.StatusOK, assemblers.BuildImageWithThumbnails(image))
}
