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
		facilitiesRoutes.GET("/user/:userId", handlers.GetAllFacilitiesForUser)

		facilitiesRoutes.POST("/", handlers.CreateFacility)
	}

	r.Run(":8080")

	return r
}
