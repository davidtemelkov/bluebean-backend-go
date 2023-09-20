package data

import (
	"context"
	"time"

	"bitbucket.org/nemetschek-systems/bluebean-service/internal/errorconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/generalconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/validator"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
)

type Space struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
	SchemaURL  string `json:"schemaURL"`
	FacilityID string `json:"facilityID"`
}

type SpaceModel struct {
	DB *dynamodb.DynamoDB
}

func ValidateSpace(v *validator.Validator, space *Space) {
	v.Check(space.Name != "", "name", errorconstants.RequiredFieldError.Error())
	v.Check(len(space.Name) >= 2, "name", errorconstants.SpaceNameMinLengthError.Error())
	v.Check(len(space.Name) < 50, "name", errorconstants.SpaceNameMaxLengthError.Error())
	v.Check(space.Location != "", "location", errorconstants.RequiredFieldError.Error())
	v.Check(len(space.Location) > 5, "location", errorconstants.SpaceLocationMinLengthError.Error())
	v.Check(len(space.Location) < 100, "location", errorconstants.SpaceLocationMaxLengthError.Error())
}

func (sm SpaceModel) Insert(space *Space) (uuid.UUID, error) {
	id := uuid.New()

	item := map[string]*dynamodb.AttributeValue{
		generalconstants.PK: {
			S: aws.String(
				generalconstants.FacilityPrefix + space.FacilityID,
			),
		},
		generalconstants.SK: {
			S: aws.String(
				generalconstants.SpacePrefix + id.String(),
			),
		},
		"ID": {
			S: aws.String(id.String()),
		},
		"Name": {
			S: aws.String(space.Name),
		},
		"Location": {
			S: aws.String(space.Location),
		},
		"SchemaURL": {
			S: aws.String(space.SchemaURL),
		},
		"FacilityID": {
			S: aws.String(space.FacilityID),
		},
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(generalconstants.TableName),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := sm.DB.PutItemWithContext(ctx, input)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (sm SpaceModel) Get(spaceID, facilityID string) (*Space, error) {
	if spaceID == "" || facilityID == "" {
		return nil, errorconstants.RecordNotFoundError
	}

	pk := generalconstants.FacilityPrefix + facilityID
	sk := generalconstants.SpacePrefix + spaceID

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

	result, err := sm.DB.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, errorconstants.RecordNotFoundError
	}

	item := result.Items[0]
	space := &Space{}
	err = dynamodbattribute.UnmarshalMap(item, space)
	if err != nil {
		return nil, err
	}

	return space, nil
}
