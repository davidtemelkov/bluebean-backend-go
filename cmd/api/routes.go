package main

import (
	"github.com/gin-gonic/gin"
)

func setupRoutes() *gin.Engine {

	r := gin.Default()
	r.Run(":8080")

	return r
}
