package data

import (
	"context"
	"example/bluebean-go/internal/validator"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Facility struct {
	Name     string   `json:"name"`
	Owners   []string `json:"owners"`
	Address  string   `json:"address"`
	City     string   `json:"city"`
	Creator  string   `json:"creator"`
	ImageUrl string   `json:"imageurl"`
}
type FacilityModel struct {
	DB *dynamodb.DynamoDB
}

func ValidateFacility(v *validator.Validator, facility *Facility) {
	v.Check(facility.Name != "", "name", "must be provided")
	v.Check(len(facility.Owners) != 0, "owners", "must be provided")
	v.Check(validator.Unique(facility.Owners), "owners", "must not contain duplicate values")
	v.Check(facility.Address != "", "address", "must be provided")
	v.Check(facility.City != "", "city", "must be provided")
	v.Check(facility.Creator != "", "creator", "must be provided")
	v.Check(facility.ImageUrl != "", "imageurl", "must be provided")
}

func (fm FacilityModel) InsertFacility(facility *Facility) error {
	item := map[string]*dynamodb.AttributeValue{
		"PK": {
			S: aws.String("FACILITY#" + facility.Name),
		},
		"SK": {
			S: aws.String("FACILITY#" + facility.Name),
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
			S: aws.String(facility.ImageUrl),
		},
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("Bluebean"),
	}

	_, err := fm.DB.PutItem(input)
	return err
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

	// facility := Facility{
	// 	Name:     aws.StringValue(result.Item["Name"].S),
	// 	Address:  aws.StringValue(result.Item["Address"].S),
	// 	City:     aws.StringValue(result.Item["City"].S),
	// 	Creator:  aws.StringValue(result.Item["Creator"].S),
	// 	ImageUrl: aws.StringValue(result.Item["ImageURL"].S),
	// }

	// ownersAttribute := result.Item["Owners"]
	// if ownersAttribute != nil && ownersAttribute.SS != nil {
	// 	facility.Owners = make([]string, len(ownersAttribute.SS))
	// 	for i, owner := range ownersAttribute.SS {
	// 		facility.Owners[i] = aws.StringValue(owner)
	// 	}
	// }

	facility := &Facility{}
	err = dynamodbattribute.UnmarshalMap(result.Item, facility)
	if err != nil {
		return nil, err
	}

	return facility, nil
}
