package data

import (
	"example/bluebean-go/internal/validator"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
