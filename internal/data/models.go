package data

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// var (
// 	ErrRecordNotFound = errors.New("record not found")
// 	ErrEditConflict   = errors.New("edit conflict")
// )

type Models struct {
	Users UserModel
}

func NewModels(db *dynamodb.DynamoDB) Models {
	return Models{
		Users: UserModel{DB: db},
	}
}
