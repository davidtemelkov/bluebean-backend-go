package main

import "example/bluebean-go/internal/utils"

// type application struct {
// 	logger *log.Logger
// 	models data.Models
// 	wg     sync.WaitGroup
// }

func main() {
	utils.LoadEnv()
	SetupRoutes()
}
