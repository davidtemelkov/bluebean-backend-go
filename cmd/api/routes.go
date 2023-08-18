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
		usersRoutes.GET("/:email/facilities", app.getFacilitiesForUserHandler)
	}

	facilitiesRoutes := r.Group("/facilities")
	{
		facilitiesRoutes.POST("/", app.createFacilityHandler)
		facilitiesRoutes.GET("/:id", app.getFacilityHandler)
	}

	r.Run(":8080")

	return r
}
