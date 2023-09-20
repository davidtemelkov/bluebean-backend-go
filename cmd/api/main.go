package main

import (
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/data"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/errorconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/mailer"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/utils"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type application struct {
	models data.Models
	mailer mailer.Mailer
}

func main() {
	utils.LoadEnv()

	db, err := openDb()
	if err != nil {
		panic(errorconstants.DBConnectionError.Error())
	}

	app := &application{
		models: data.NewModels(db),
		mailer: mailer.New(utils.GetSMTPHost(), utils.GetSMTPPort(), utils.GetSMTPUsername(), utils.GetSMTPPassword(), utils.GetSMTPSender()),
	}

	app.setupRoutes()
}

func openDb() (*dynamodb.DynamoDB, error) {
	awsAccessKeyID := utils.GetAWSAccessKey()
	awsSecretAccessKey := utils.GetAWSSecretKey()

	sess, err := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),

		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})

	if err != nil {
		return nil, err
	}

	return dynamodb.New(sess), nil
}
