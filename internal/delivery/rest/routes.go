package rest

import (
	"github.com/gin-gonic/gin"
	"image-resizing-shared/internal/delivery/rest/handlers"
	"image-resizing-shared/pkg/di"
	"net/http"
	"os"
)

func RegisterRoutes(router *gin.Engine, dependencies *di.Dependencies) {
	imageHandler := handlers.NewImageHandler(dependencies.ImageService)

	router.GET("/ping", imageHandler.Ping)
	router.POST("/image/upload", UploadAuthMiddleware(), imageHandler.UploadImage)
	router.POST("/image/upload/binary", UploadAuthMiddleware(), imageHandler.UploadImageBinary)
	router.GET("/image/:id", imageHandler.GetImage)

	router.StaticFS("/images", gin.Dir(dependencies.UploadsDir, false))
}

func UploadAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uploadToken := os.Getenv("REST_UPLOAD_TOKEN")
		if uploadToken == "" {
			c.Next()
			return
		}

		if c.GetHeader("X-API-Key") != uploadToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		c.Next()
	}
}
