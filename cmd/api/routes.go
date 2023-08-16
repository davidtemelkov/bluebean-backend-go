package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) setupRoutes() *gin.Engine {

	r := gin.Default()

	usersRoutes := r.Group("/users")
	{
		usersRoutes.POST("/", app.regiserUserHandler)
	}

	r.Run(":8080")

	return r
}
