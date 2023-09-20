package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (app *application) setupRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           0,
	}))

	usersRoutes := r.Group("/users")
	{
		usersRoutes.POST("/register", app.registerUserHandler)
		usersRoutes.POST("/login", app.loginUserHandler)
		usersRoutes.Use(app.authenticate())
		usersRoutes.GET("/:email/facilities", app.getAllFacilitiesForUserHandler)
	}

	facilitiesRoutes := r.Group("/facilities")
	{
		facilitiesRoutes.Use(app.authenticate())
		facilitiesRoutes.POST("/", app.createFacilityHandler)
		facilitiesRoutes.GET("/:facilityID", app.getFacilityHandler)
		facilitiesRoutes.POST("/users", app.addUserToFacilityHandler)
		facilitiesRoutes.DELETE("/:facilityID/user/:email", app.removeUserFromFacilityHandler)
		facilitiesRoutes.GET("/:facilityID/users", app.getAllUsersForFacility)
		facilitiesRoutes.GET("/:facilityID/spaces", app.getAllSpacesForFacility)
		facilitiesRoutes.PATCH("/assets/add", app.addAssetToFacilityHandler)
		facilitiesRoutes.PATCH("/assets/remove", app.removeAssetFromFacilityHandler)
	}

	spacesRoutes := r.Group("/spaces")
	{
		spacesRoutes.Use(app.authenticate())
		spacesRoutes.POST("/", app.createSpaceHandler)
		spacesRoutes.GET("/:spaceID/facility/:facilityID", app.getSpaceHandler)
	}

	punchesRoutes := r.Group("/punches")
	{
		punchesRoutes.Use(app.authenticate())
		punchesRoutes.POST("/:facilityID", app.createPunchHandler)
		punchesRoutes.GET("/:punchID/facility/:facilityID/space/:spaceID", app.getPunchHandler)
		punchesRoutes.GET("/facility/:facilityID/space/:spaceID", app.getAllPunchesForSpaceHandler)
		punchesRoutes.GET("/facility/:facilityID", app.getAllPunchesForFacilityHandler)
		punchesRoutes.PUT("/:facilityID", app.editPunchHandler)
		punchesRoutes.DELETE("/:punchID/facility/:facilityID/space/:spaceID", app.deletePunchHandler)
	}

	commentsRoutes := r.Group("/comments")
	{
		commentsRoutes.Use(app.authenticate())
		commentsRoutes.POST("/", app.createCommentHandler)
		commentsRoutes.GET("/:facilityID/space/:spaceID/punch/:punchID", app.getAllCommentsForPunchHandler)
	}

	r.Run(":8080")

	return r
}
