package router

import (

	"github.com/gin-gonic/gin"
	"harbor-tools/harbor-tools/controllers"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	apiV1 := router.Group("v1")
	{
		apiV1.GET("/test", controllers.Test)
		apiV1.POST("/tags", )
		//v1.GET()
	}
	return router
}
