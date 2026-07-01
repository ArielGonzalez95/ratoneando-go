package routes

import (
	"github.com/gin-gonic/gin"

	"ratoneando/controllers"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/", controllers.NormalizedScraper)
	router.GET("/raw", controllers.NormalizedScraper)

	router.POST("/lista", controllers.Lista)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
