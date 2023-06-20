package routes

import (
	"example/bluebean-go/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	r := gin.Default()

	facilitiesRoutes := r.Group("/facilities")
	{
		facilitiesRoutes.GET("/", handlers.GetAllFacilities)
		facilitiesRoutes.GET("/:facilityId", handlers.GetOneFacility)
		//facilitiesRoutes.GET("/:facilityId/:userId", handlers.GetAllForUser)
	}

	r.Run(":8080")

	return r
}
