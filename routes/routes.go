package routes

import (
	"github.com/99-66/go-gin-token-server/controllers"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		token := api.Group("/token")
		{
			token.POST("/create",  controllers.CreateToken)
			token.GET("/verify", controllers.VerifyToken)
			token.POST("/refresh", controllers.RefreshToken)
		}
	}

	return r
}