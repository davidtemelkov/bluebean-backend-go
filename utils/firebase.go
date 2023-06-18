package utils

import (
	"context"
	"encoding/base64"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type FireBaseStorage struct {
	Bucket string
}

func NewFireBaseStorage(bucket string) *FireBaseStorage {
	return &FireBaseStorage{
		Bucket: bucket,
	}
}

func UploadFile(photo64, fileFolder, fileName string) (string, error) {
	bucketName := GetFirebaseBucketName()

	fb := NewFireBaseStorage(bucketName)
	ctx := context.Background()
	opt := option.WithCredentialsFile("utils/serviceAccountKey.json")

	client, err := storage.NewClient(ctx, opt)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Storage client: %v", err)
	}

	photoData, err := base64.StdEncoding.DecodeString(photo64)
	if err != nil {
		return "", err
	}

	filePath := fileFolder + "/" + fileName
	wc := client.Bucket(fb.Bucket).Object(filePath).NewWriter(ctx)
	wc.ContentType = "image/jpeg" // Set the appropriate content type if needed
	if _, err := wc.Write(photoData); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	firebaseUrl := GetFirebaseUrl()
	url := "https://firebasestorage.googleapis.com" + firebaseUrl + fileFolder + "%2F" + fileName + "?alt=media"

	return url, nil
}
