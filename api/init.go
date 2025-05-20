package api

import (
	"image-captioning-go/api/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitServer() {
	
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5500"}, 
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))
	r.Static("/images", "./images")

	
	apiRoute := r.Group("/api")
	{
		uh := handlers.UploadHandler{}
		apiRoute.POST("/upload", uh.Upload)
	}
	r.Run(":8080")
}