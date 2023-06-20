package routes

import (
	"example/bluebean-go/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	r := gin.Default()

	facilitiesRoutes := r.Group("/facilities")
	{
		facilitiesRoutes.GET("/", handlers.GetOneFacility)
	}

	r.Run(":8080")

	return r
}
