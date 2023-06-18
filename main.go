package main

import (
	"example/bluebean-go/database"
	"example/bluebean-go/utils"
	"fmt"
)

func main() {

	utils.LoadEnv()
	database.SetupDB()

	test := utils.GetTest()
	url, err := utils.UploadFile(test, "facilities", "CoverPhoto")
	if err == nil {
		fmt.Println(url)
	} else {
		fmt.Println(err)
	}

	//routes.SetupRoutes()
}
