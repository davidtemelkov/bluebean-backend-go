package data

import (
	"context"
	"time"

	"bitbucket.org/nemetschek-systems/bluebean-service/internal/errorconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/generalconstants"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type UserFacility struct {
	Username         string `json:"username"`
	UserEmail        string `json:"userEmail"`
	UserRole         string `json:"userRole"`
	UserAddedOn      string `json:"userAddedOn"`
	FacilityID       string `json:"facilityID"`
	FacilityName     string `json:"facilityName"`
	FacilityAddress  string `json:"facilityAddress"`
	FacilityCity     string `json:"facilityCity"`
	FacilityImageURL string `json:"facilityImageURL"`
	GSI1PK           string `json:"GSI1PK"`
	GSI1SK           string `json:"GSI1SK"`
}

type UserFacilityModel struct {
	DB *dynamodb.DynamoDB
}

func (ufm UserFacilityModel) Get(userEmail string, facilityID string) (*UserFacility, error) {
	if userEmail == "" || facilityID == "" {
		return nil, errorconstants.RecordNotFoundError
	}

	pk := generalconstants.UserPrefix + userEmail
	sk := generalconstants.FacilityPrefix + facilityID

	keyCondition := expression.Key(generalconstants.PK).Equal(expression.Value(pk)).
		And(expression.Key(generalconstants.SK).Equal(expression.Value(sk)))

	builder, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		return nil, err
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(generalconstants.TableName),
		KeyConditionExpression:    builder.KeyCondition(),
		ExpressionAttributeNames:  builder.Names(),
		ExpressionAttributeValues: builder.Values(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := ufm.DB.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, errorconstants.RecordNotFoundError
	}

	item := result.Items[0]
	userFacility := &UserFacility{}
	err = dynamodbattribute.UnmarshalMap(item, userFacility)
	if err != nil {
		return nil, err
	}

	return userFacility, nil
}

func (ufm UserFacilityModel) Insert(user *User, facility *Facility) error {
	item := map[string]*dynamodb.AttributeValue{
		generalconstants.PK: {
			S: aws.String(
				generalconstants.UserPrefix + user.Email,
			),
		},
		generalconstants.SK: {
			S: aws.String(
				generalconstants.FacilityPrefix + facility.ID,
			),
		},
		"FacilityID":       {S: aws.String(facility.ID)},
		"FacilityName":     {S: aws.String(facility.Name)},
		"FacilityCity":     {S: aws.String(facility.City)},
		"FacilityAddress":  {S: aws.String(facility.Address)},
		"FacilityImageURL": {S: aws.String(facility.ImageURL)},
		"UserEmail":        {S: aws.String(user.Email)},
		"UserName":         {S: aws.String(user.Name)},
		"UserRole":         {S: aws.String(user.Role)},
		"UserAddedOn":      {S: aws.String(time.Now().UTC().Format(time.RFC3339))},
		generalconstants.GSI1PK: {
			S: aws.String(
				generalconstants.FacilityPrefix + facility.ID,
			),
		},
		generalconstants.GSI1SK: {
			S: aws.String(
				generalconstants.UserPrefix + user.Email,
			),
		},
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(generalconstants.TableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := ufm.DB.PutItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
