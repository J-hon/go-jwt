package routes

import (
	controllers "github.com/J-hon/controllers"
	"github.com/gin-gonic/gin"
)

func Auth(inc *gin.Engine) {
	inc.POST("/signup", controllers.SignUp)
	inc.POST("/signin", controllers.SignIn)
}
