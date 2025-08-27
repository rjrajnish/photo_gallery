package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rjrajnish/photo_gallery/backend/cmd/api/middleware"
	"github.com/rjrajnish/photo_gallery/backend/internal/handlers"
)

func SetupRouter(r *gin.Engine) {
	auth := handlers.NewAuthHandler()
	folders := handlers.NewFolderHandler()
	photos := handlers.NewPhotoHandler()

	// public routes
	r.POST("/api/auth/register", auth.Register)
	r.POST("/api/auth/login", auth.Login)

	// secured routes
	sec := r.Group("/api")
	sec.Use(middleware.AuthRequired())

	sec.GET("/folders", folders.List)
	sec.POST("/folders", folders.Create)
	sec.DELETE("/folders", folders.Delete)

	sec.GET("/photos", photos.List)
	sec.POST("/photos/upload", photos.Upload)
	sec.POST("/photos/delete", photos.Delete)
	st := r.Group("/api")
	st.GET("/photos/:id/stream", photos.Stream)
}
