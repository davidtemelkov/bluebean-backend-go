package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"bitbucket.org/nemetschek-systems/bluebean-service/internal/errorconstants"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type FireBaseStorage struct {
	Bucket string
}

// Folder constants
var (
	FacilitiesFolder = "facility_images"
	SpacesFolder     = "spaces_images"
)
var expectedPrefixes = map[string]string{
	"image/jpeg;base64,": "image/jpeg",
	"image/png;base64,":  "image/png",
}

func NewFireBaseStorage(bucket string) *FireBaseStorage {
	return &FireBaseStorage{
		Bucket: bucket,
	}
}

func ValidateAndExtractContentType(photo64 string) (string, string, error) {
	for prefix, contentType := range expectedPrefixes {
		if strings.HasPrefix(photo64, prefix) {
			return contentType, prefix, nil
		}
	}

	return "", "", errorconstants.InvalidBase64ImagePrefixError
}

func UploadFile(photo64, fileFolder, fileName string) (string, error) {
	bucketName := GetFirebaseBucketName()

	fb := NewFireBaseStorage(bucketName)
	ctx := context.Background()
	opt := option.WithCredentialsFile("internal/utils/serviceAccountKey.json")

	client, err := storage.NewClient(ctx, opt)
	if err != nil {
		panic(errorconstants.FirebaseClientError.Error())
	}

	contentType, prefix, err := ValidateAndExtractContentType(photo64)
	if err != nil {
		return "", err
	}

	photo64 = strings.TrimPrefix(photo64, prefix)

	photoData, err := base64.StdEncoding.DecodeString(photo64)
	if err != nil {
		return "", err
	}

	fileName = strings.ReplaceAll(fileName, " ", "")
	filePath := fmt.Sprintf("%s/%s", fileFolder, fileName)
	wc := client.Bucket(fb.Bucket).Object(filePath).NewWriter(ctx)

	wc.ContentType = contentType

	if _, err := wc.Write(photoData); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	url, err := generateFirebaseUrl(fileFolder, fileName)
	if err != nil {
		return "", err
	}

	return url, nil
}

func generateFirebaseUrl(fileFolder, fileName string) (string, error) {
	baseUrl := GetFirebaseUrl()

	if fileFolder == "" {
		return "", errorconstants.FileFolderEmptyError
	}
	if fileName == "" {
		return "", errorconstants.FileNameEmptyError
	}

	url := fmt.Sprintf("%s%s%%2F%s?alt=media", baseUrl, fileFolder, fileName)

	return url, nil
}
