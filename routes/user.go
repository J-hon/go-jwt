package routes

import (
	"github.com/J-hon/controllers"
	"github.com/J-hon/go-jwt/middleware"
	"github.com/gin-gonic/gin"
)

func User(inc *gin.Engine) {
	inc.Use(middleware.Authenticate)

	inc.GET("/users", controllers.GetUsers)
	inc.GET("/users/:id", controllers.GetUser)
}
