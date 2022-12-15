package routes

import (
	"github.com/J-hon/go-jwt/controllers"
	"github.com/gin-gonic/gin"
)

func Auth(inc *gin.Engine) {
	inc.POST("/signup", controllers.SignUp())
	inc.POST("/signin", controllers.SignIn())
}
