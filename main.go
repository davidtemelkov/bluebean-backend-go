package main

import (
	"log"

	"github.com/joho/godotenv"
	"gorm.io/gen"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	setupDB()

	conf := gen.Config{
		OutPath: "./query",
		Mode:    gen.WithDefaultQuery, // generate mode
	}
	g := gen.NewGenerator(conf)

	g.UseDB(db) // reuse your gorm db
	g.ApplyBasic(
		// Generate structs from all tables of current database
		g.GenerateAllTable()...,
	)
	g.Execute()
}
