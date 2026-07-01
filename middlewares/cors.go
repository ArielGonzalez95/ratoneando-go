package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"ratoneando/config"
)

func CORS(router *gin.Engine) {
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: false,
	}
	if config.ENV == "release" && config.WEB_URL != "" {
		corsConfig.AllowOrigins = []string{config.WEB_URL}
	}
	router.Use(cors.New(corsConfig))
}
