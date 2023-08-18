package data

import (
	"context"
	"example/bluebean-go/internal/validator"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type Facility struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Owners   []string `json:"owners"`
	Address  string   `json:"address"`
	City     string   `json:"city"`
	Creator  string   `json:"creator"`
	ImageURL string   `json:"imageurl"`
}
type FacilityModel struct {
	DB    *dynamodb.DynamoDB
	Users UserModel
}

func ValidateFacility(v *validator.Validator, facility *Facility) {
	v.Check(facility.Name != "", "name", "must be provided")
	v.Check(len(facility.Owners) != 0, "owners", "must be provided")
	v.Check(validator.Unique(facility.Owners), "owners", "must not contain duplicate values")
	v.Check(facility.Address != "", "address", "must be provided")
	v.Check(facility.City != "", "city", "must be provided")
	v.Check(facility.Creator != "", "creator", "must be provided")
	v.Check(facility.ImageURL != "", "imageurl", "must be provided")
}

func (fm FacilityModel) InsertFacility(facility *Facility) (uuid.UUID, error) {
	id := uuid.New()

	item := map[string]*dynamodb.AttributeValue{
		"PK": {
			S: aws.String("FACILITY#" + id.String()),
		},
		"SK": {
			S: aws.String("FACILITY#" + id.String()),
		},
		"ID": {
			S: aws.String(id.String()),
		},
		"Name": {
			S: aws.String(facility.Name),
		},
		"Owners": {
			SS: aws.StringSlice(facility.Owners),
		},
		"Address": {
			S: aws.String(facility.Address),
		},
		"City": {
			S: aws.String(facility.City),
		},
		"Creator": {
			S: aws.String(facility.Creator),
		},
		"ImageURL": {
			S: aws.String(facility.ImageURL),
		},
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("Bluebean"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := fm.DB.PutItemWithContext(ctx, input)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (fm FacilityModel) Get(id string) (*Facility, error) {
	if id == "" {
		return nil, ErrRecordNotFound
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String("Bluebean"),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String("FACILITY#" + id),
			},
			"SK": {
				S: aws.String("FACILITY#" + id),
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := fm.DB.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.Item) == 0 {
		return nil, ErrRecordNotFound
	}

	facility := &Facility{}
	err = dynamodbattribute.UnmarshalMap(result.Item, facility)
	if err != nil {
		return nil, err
	}

	return facility, nil
}

func (fm FacilityModel) AddUserToFacility(userEmail string, id string, um UserModel) error {
	facility, err := fm.Get(id)
	if err != nil {
		return ErrRecordNotFound
	}

	user, err := um.GetByEmail(userEmail)
	if err != nil {
		return ErrRecordNotFound
	}

	item := map[string]*dynamodb.AttributeValue{
		"PK": {
			S: aws.String("USER#" + userEmail),
		},
		"SK": {
			S: aws.String("FACILITY#" + id),
		},
		"FacilityID":       {S: aws.String(id)},
		"FacilityName":     {S: aws.String(facility.Name)},
		"FacilityAddress":  {S: aws.String(facility.Address)},
		"FacilityImageURL": {S: aws.String(facility.ImageURL)},
		"UserEmail":        {S: aws.String(userEmail)},
		"UserName":         {S: aws.String(user.Name)},
		"UserRole":         {S: aws.String(user.Role)},
		"UserAddedOn":      {S: aws.String(time.Now().String())},
		"GSI1PK": {
			S: aws.String("FACILITY#" + id),
		},
		"GSI1SK": {
			S: aws.String("USER#" + userEmail),
		},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Bluebean"),
		Item:      item,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = fm.DB.PutItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (fm FacilityModel) GetAllUsersForFacility(id string) ([]User, error) {
	keyConditionExpression := "GSI1PK = :gsi1pk AND begins_with(GSI1SK, :gsi1skPrefix)"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":gsi1pk": {
			S: aws.String("FACILITY#" + id),
		},
		":gsi1skPrefix": {
			S: aws.String("USER#"),
		},
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String("Bluebean"),
		IndexName:                 aws.String("GSI1"),
		KeyConditionExpression:    aws.String(keyConditionExpression),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := fm.DB.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)

	for _, item := range result.Items {
		user := User{
			Email:   *item["UserEmail"].S,
			Name:    *item["UserName"].S,
			Role:    *item["UserRole"].S,
			AddedOn: *item["UserAddedOn"].S,
		}
		users = append(users, user)
	}

	return users, nil
}
