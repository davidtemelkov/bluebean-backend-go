package data

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users UserModel
	Facilities FacilityModel
}

func NewModels(db *dynamodb.DynamoDB) Models {
	return Models{
		Users: UserModel{DB: db},
		Facilities: FacilityModel{DB: db},
	}
}
