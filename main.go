package main

import (
	"example/bluebean-go/database"
	"example/bluebean-go/utils"
)

func main() {
	utils.LoadEnv()
	database.SetupDB()
}
