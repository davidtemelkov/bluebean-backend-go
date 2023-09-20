package data

import (
	"context"
	"time"

	"bitbucket.org/nemetschek-systems/bluebean-service/internal/errorconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/generalconstants"

	"bitbucket.org/nemetschek-systems/bluebean-service/internal/validator"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
)

type Comment struct {
	ID           string `json:"id"`
	PunchID      string `json:"punchID"`
	SpaceID      string `json:"spaceID"`
	FacilityID   string `json:"facilityID"`
	Text         string `json:"text"`
	CreatedOn    string `json:"createdOn"`
	CreatorEmail string `json:"creatorEmail"`
	CreatorName  string `json:"creatorName"`
}

type CommentModel struct {
	DB *dynamodb.DynamoDB
}

func ValidateComment(v *validator.Validator, comment *Comment) {
	v.Check(comment.PunchID != "", "punchId", errorconstants.RequiredFieldError.Error())
	v.Check(comment.SpaceID != "", "spaceId", errorconstants.RequiredFieldError.Error())
	v.Check(comment.FacilityID != "", "facilityId", errorconstants.RequiredFieldError.Error())
	v.Check(comment.Text != "", "text", errorconstants.RequiredFieldError.Error())
	v.Check(len(comment.Text) > 5, "text", errorconstants.CommentTextMinLengthError.Error())
	v.Check(len(comment.Text) < 500, "text", errorconstants.CommentTextMaxLengthError.Error())
}

func (cm CommentModel) Insert(comment *Comment) (uuid.UUID, error) {
	id := uuid.New()

	item := map[string]*dynamodb.AttributeValue{
		generalconstants.PK: {
			S: aws.String(
				generalconstants.FacilityPrefix + comment.FacilityID +
					generalconstants.SpacePrefix + comment.SpaceID,
			),
		},
		generalconstants.SK: {
			S: aws.String(
				generalconstants.PunchPrefix + comment.PunchID +
					generalconstants.CommentPrefix + id.String(),
			),
		},
		"ID": {
			S: aws.String(id.String()),
		},
		"PunchID": {
			S: aws.String(comment.FacilityID),
		},
		"SpaceID": {
			S: aws.String(comment.SpaceID),
		},
		"FacilityID": {
			S: aws.String(comment.FacilityID),
		},
		"Text": {
			S: aws.String(comment.Text),
		},
		"CreatedOn": {
			S: aws.String(comment.CreatedOn),
		},
		"CreatorEmail": {
			S: aws.String(comment.CreatorEmail),
		},
		"CreatorName": {
			S: aws.String(comment.CreatorName),
		},
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(generalconstants.TableName),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := cm.DB.PutItemWithContext(ctx, input)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (cm CommentModel) GetAllCommentsForPunch(punchID, spaceID, facilityID string) ([]Comment, error) {
	keyCondition := expression.Key(generalconstants.PK).
		Equal(
			expression.Value(
				generalconstants.FacilityPrefix + facilityID +
					generalconstants.SpacePrefix + spaceID,
			),
		).And(
		expression.Key(generalconstants.SK).
			BeginsWith(
				generalconstants.PunchPrefix + punchID +
					generalconstants.CommentPrefix,
			),
	)

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

	result, err := cm.DB.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	comments := make([]Comment, 0)

	for _, item := range result.Items {
		comment := Comment{
			ID:           *item["ID"].S,
			PunchID:      *item["PunchID"].S,
			SpaceID:      *item["SpaceID"].S,
			FacilityID:   *item["FacilityID"].S,
			Text:         *item["Text"].S,
			CreatedOn:    *item["CreatedOn"].S,
			CreatorEmail: *item["CreatorEmail"].S,
			CreatorName:  *item["CreatorName"].S,
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
