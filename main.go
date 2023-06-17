package main

import (
	"example/bluebean-go/routes"
	"example/bluebean-go/database"
	"example/bluebean-go/utils"
)

func main() {	
  
	utils.LoadEnv()
	database.SetupDB() 
  
  router := routes.SetupRoutes()
	router.Run(":8080");
}
