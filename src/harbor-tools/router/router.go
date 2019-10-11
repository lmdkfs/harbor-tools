package router

import (

	"github.com/gin-gonic/gin"
	"harbor-tools/harbor-tools/controllers"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("v1")
	{
		v1.GET("/test", controllers.Test)
	}
	return router
}
