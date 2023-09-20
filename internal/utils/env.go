package utils

import (
	"os"
	"strconv"

	"bitbucket.org/nemetschek-systems/bluebean-service/internal/errorconstants"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(errorconstants.LoadingEnvFileError.Error())
	}
}

func GetFirebaseUrl() string {
	firebaseUrl := os.Getenv("FIREBASE_URL")

	if firebaseUrl == "" {
		panic(errorconstants.FirebaseURLError.Error())
	}

	return firebaseUrl
}

func GetFirebaseBucketName() string {
	firebaseBucketName := os.Getenv("FIREBASE_BUCKET_NAME")

	if firebaseBucketName == "" {
		panic(errorconstants.FirebaseBucketNameError.Error())
	}

	return firebaseBucketName
}

func GetAWSAccessKey() string {
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")

	if awsAccessKey == "" {
		panic(errorconstants.AWSAccessKeyError.Error())
	}

	return awsAccessKey
}

func GetAWSSecretKey() string {
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")

	if awsSecretKey == "" {
		panic(errorconstants.AWSSecretKeyError.Error())
	}

	return awsSecretKey
}

func GetJWTPrivateKey() []byte {
	jwtPrivateKey := os.Getenv("JWT_PRIVATE_KEY")

	if jwtPrivateKey == "" {
		panic(errorconstants.JWTPrivateKeyError.Error())
	}

	return []byte(jwtPrivateKey)
}

func GetSMTPHost() string {
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		panic(errorconstants.SMTPHostError.Error())
	}
	return smtpHost
}

func GetSMTPPort() int {
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		panic(errorconstants.SMTPPortError.Error())
	}
	return smtpPort
}

func GetSMTPUsername() string {
	smtpUsername := os.Getenv("SMTP_USERNAME")
	if smtpUsername == "" {
		panic(errorconstants.SMTPUsernameError.Error())
	}
	return smtpUsername
}

func GetSMTPPassword() string {
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
		panic(errorconstants.SMTPPasswordError.Error())
	}
	return smtpPassword
}

func GetSMTPSender() string {
	smtpSender := os.Getenv("SMTP_SENDER")
	if smtpSender == "" {
		panic(errorconstants.SMTPSenderError.Error())
	}
	return smtpSender
}

func GetWebAppBaseUrl() string {
	webAppBaseUrl := os.Getenv("WEB_APP_BASE_URL")
	if webAppBaseUrl == "" {
		panic(errorconstants.WebAppBaseUrlError.Error())
	}
	return webAppBaseUrl
}
