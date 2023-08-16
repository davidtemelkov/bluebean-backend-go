package data

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// var (
// 	ErrRecordNotFound = errors.New("record not found")
// 	ErrEditConflict   = errors.New("edit conflict")
// )

type Models struct {
	Facilities FacilityModel
}

func NewModels(db *dynamodb.DynamoDB) Models {
	return Models{
		Facilities: FacilityModel{DB: db},
	}
}
