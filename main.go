package main

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/tonacobr/auth-service/controller"
	_ "github.com/tonacobr/auth-service/docs"
	"net/http"
)

// @title          Golang AuthService
// @version        1.0.0
// @description    JWT service which will generate and verify jwt tokens
// @termsOfService http://swagger.io/terms/

// @contact.name  Team Eaglets
// @contact.url   https://team-eaglets-fwt.io/support
// @contact.email eagles.v3@team-eagles.io

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     localhost:8081
// @BasePath /api/v1
func main() {
	router := gin.Default()

	// init default controller
	c := controller.NewController()

	// health-check
	router.GET("/ping", healthCheck)

	// api endpoints (v1)
	v1 := router.Group("/api/v1")
	{
		// jwt controller
		jwt := v1.Group("/jwt")
		{
			jwt.POST("", c.GenerateToken)
			jwt.GET("", c.GetPublicKey)
			jwt.GET(":token", c.ValidateToken)
		}
	}

	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// run application - TODO: research best practices for launching gin server
	err := router.Run(":8081")
	if err != nil {
		panic(err)
	}
}

func healthCheck(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
