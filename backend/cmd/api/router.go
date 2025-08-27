package main

import (
	"github.com/gin-gonic/gin"
	"github.com/Rjrajnish99/photo-gallery/backend/cmd/api/middleware"
	"github.com/Rjrajnish99/photo-gallery/backend/internal/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	auth := handlers.NewAuthHandler()
	folders := handlers.NewFolderHandler()
	photos := handlers.NewPhotoHandler()

	r.POST("/api/auth/register", auth.Register)
	r.POST("/api/auth/login", auth.Login)

	sec := r.Group("/api")
	sec.Use(middleware.AuthRequired())

	sec.GET("/folders", folders.List)
	sec.POST("/folders", folders.Create)
	sec.DELETE("/folders", folders.Delete)

	sec.GET("/photos", photos.List)
	sec.POST("/photos/upload", photos.Upload)
	sec.POST("/photos/delete", photos.Delete)
	sec.GET("/photos/:id/stream", photos.Stream)

	return r
}
