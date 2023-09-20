package data

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Models struct {
	Users          UserModel
	Facilities     FacilityModel
	UserFacilities UserFacilityModel
	Spaces         SpaceModel
	Punches        PunchModel
	Comments       CommentModel
}

func NewModels(db *dynamodb.DynamoDB) Models {
	return Models{
		Users:          UserModel{DB: db},
		Facilities:     FacilityModel{DB: db},
		UserFacilities: UserFacilityModel{DB: db},
		Spaces:         SpaceModel{DB: db},
		Punches:        PunchModel{DB: db},
		Comments:       CommentModel{DB: db},
	}
}
