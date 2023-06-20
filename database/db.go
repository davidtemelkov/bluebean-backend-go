package database

import (
	"example/bluebean-go/utils"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var Db *gorm.DB

func SetupDB() {
	dbConnectionString := utils.GetConnectionString()

	var err error
	Db, err = gorm.Open(sqlserver.Open(dbConnectionString), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}
