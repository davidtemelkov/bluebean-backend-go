package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) setupRoutes() *gin.Engine {

	r := gin.Default()

	facilitiesRoutes := r.Group("/facilities")
	{
		facilitiesRoutes.POST("/", app.createFacilityHandler)
	}

	r.Run(":8080")

	return r
}
