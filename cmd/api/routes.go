package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) setupRoutes() *gin.Engine {
	r := gin.Default()

	usersRoutes := r.Group("/users")
	{
		usersRoutes.POST("/", app.regiserUserHandler)
		usersRoutes.GET("/:email", app.getUserByEmailHandler)
	}

	facilitiesRoutes := r.Group("/facilities")
	{
		facilitiesRoutes.POST("/", app.createFacilityHandler)
		facilitiesRoutes.GET("/:facilityName", app.getFacilityHandler)
		facilitiesRoutes.POST("/:facilityName/user/:email", app.addUserToFacilityHandler)
		facilitiesRoutes.GET("/:facilityName/users", app.getAllUsersForFacility)
	}

	r.Run(":8080")

	return r
}
