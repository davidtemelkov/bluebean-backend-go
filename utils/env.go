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

func GetFirebaseUrl() string {
	firebaseUrl := os.Getenv("FIREBASE_URL")

	if firebaseUrl == "" {
		log.Fatal("FIREBASE_URL environment variable is not set")
	}

	return firebaseUrl
}

func GetFirebaseBucketName() string {
	firebaseBucketName := os.Getenv("FIREBASE_BUCKET_NAME")

	if firebaseBucketName == "" {
		log.Fatal("FIREBASE_URL environment variable is not set")
	}

	return firebaseBucketName
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
