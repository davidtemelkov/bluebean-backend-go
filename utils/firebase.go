package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func UploadFile(photo64 string, fileFolder string, fileName string) (string, error) {
	opt := option.WithCredentialsFile("utils/serviceAccountKey.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		log.Fatalf("Error initializing Firebase Storage client: %v", err)
	}

	fileBytes, err := base64.StdEncoding.DecodeString(photo64)
	if err != nil {
		return "", fmt.Errorf("error decoding file: %v", err)
	}

	filePath := fileFolder + "/" + fileName

	bucket, err := client.Bucket("bluebean-4d5ab.appspot.com")
	if err != nil {
		return "", fmt.Errorf("error getting bucket handle: %v", err)
	}

	writer := bucket.Object(filePath).NewWriter(context.Background())
	if _, err := writer.Write(fileBytes); err != nil {
		return "", fmt.Errorf("error uploading file to Firebase Storage: %v", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("error closing writer: %v", err)
	}
	attrs, err := bucket.Object(filePath).Attrs(context.Background())
	if err != nil {
		return "", fmt.Errorf("error retrieving file attributes: %v", err)
	}
	downloadURL := attrs.MediaLink

	return downloadURL, nil
}
