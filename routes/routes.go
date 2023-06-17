package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	r := gin.Default()
	r.Run(":8080")

	return r
}