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
