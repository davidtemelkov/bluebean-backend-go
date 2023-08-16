package main

import (
	"example/bluebean-go/internal/data"
	"example/bluebean-go/internal/utils"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type application struct {
	// logger *log.Logger
	// // models data.Models
	// wg sync.WaitGroup
	models data.Models
}

func main() {
	utils.LoadEnv()

	db, err := openDb()
	if err != nil {
		panic("Db connection error")
	}

	// app := &application{
	// 	models: data.NewModels(dbClient),
	// }

	input := &dynamodb.PutItemInput{

		TableName: aws.String("Bluebean"),

		Item: map[string]*dynamodb.AttributeValue{

			// "PK": {

			// 	S: aws.String("proba"),
			// },

			// "SK": {

			// 	S: aws.String("bountyhunter1"),
			// },

			// "Proba": {
			// 	S: aws.String("probaNaProbata"),
			// },
		},

		// Add a condition expression to check if SessionToken attribute doesn't exist

		ConditionExpression: aws.String("attribute_not_exists(SessionToken)"),
	}

	// Put the item into DynamoDB table

	_, err = db.PutItem(input)

	if err != nil {

		fmt.Println("Error putting item into DynamoDB:", err)

		return

	}

	fmt.Println("Item successfully put into DynamoDB")

	setupRoutes()

}

func openDb() (*dynamodb.DynamoDB, error) {
	awsAccessKeyID := utils.GetAWSAccessKey()
	awsSecretAccessKey := utils.GetAWSSecretKey()

	sess, err := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"), // Change this to your preferred AWS region

		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})

	if err != nil {
		fmt.Println("Error creating AWS session:", err)

		return nil, err
	}

	return dynamodb.New(sess), nil
}