package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetConnectionString() string {
	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")

	if dbConnectionString == "" {
		log.Fatal("DB_CONNECTION_STRING environment variable is not set")
	}

	return dbConnectionString
}

func GetTest() string {
	test := os.Getenv("Test")

	if test == "" {
		log.Fatal("DB_CONNECTION_STRING environment variable is not set")
	}

	return test
}

func GetJWTKey() string {
	jwtKey := os.Getenv("JWT_KEY")

	if jwtKey == "" {
		log.Fatal("JWT_KEY environment variable is not set")
	}

	return jwtKey

}

func GetJWTIssuer() string {
	jwtIssuer := os.Getenv("JWT_ISSUER")

	if jwtIssuer == "" {
		log.Fatal("JWT_ISSUER environment variable is not set")
	}

	return jwtIssuer

}

func GetJWTAudience() string {
	jwtAudience := os.Getenv("JWT_AUDIENCE")

	if jwtAudience == "" {
		log.Fatal("JWT_AUDIENCE environment variable is not set")
	}

	return jwtAudience

}
