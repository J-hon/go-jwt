package main

import (
	"os"

	"github.com/J-hon/go-jwt/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = ":9010"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.Auth(router)
	routes.User(router)

	router.GET("api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	router.GET("api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(":" + port)
}
