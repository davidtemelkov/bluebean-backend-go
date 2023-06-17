package main

import (
	"log"
	"os"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var db *gorm.DB

func setupDB() {
	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")

	if dbConnectionString == "" {
		log.Fatal("DB_CONNECTION_STRING environment variable is not set")
	}

	var err error
	db, err = gorm.Open(sqlserver.Open(dbConnectionString), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}
