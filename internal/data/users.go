package data

import (
	"context"
	"errors"
	"example/bluebean-go/internal/validator"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password password `json:"-"`
	Role     string   `json:"role"`
	AddedOn  string   `json:"addedOn,omitempty"`
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 6, "password", "must be at least 6 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

type UserModel struct {
	DB *dynamodb.DynamoDB
}

func (um UserModel) Insert(user *User) error {
	item := map[string]*dynamodb.AttributeValue{
		"PK": {
			S: aws.String("USER#" + user.Email),
		},
		"SK": {
			S: aws.String("USER#" + user.Email),
		},
		"Name": {
			S: aws.String(user.Name),
		},
		"Email": {
			S: aws.String(user.Email),
		},
		"HashedPassword": {
			S: aws.String(string(user.Password.hash)),
		},
		"Role": {
			S: aws.String(user.Role),
		},
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String("Bluebean"),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := um.DB.PutItemWithContext(ctx, input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (um UserModel) GetByEmail(email string) (*User, error) {
	key := map[string]*dynamodb.AttributeValue{
		"PK": {
			S: aws.String("USER#" + email),
		},
		"SK": {
			S: aws.String("USER#" + email),
		},
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String("Bluebean"),
		Key:       key,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := um.DB.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.Item) == 0 {
		return nil, nil // User not found
	}

	user := &User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (um UserModel) GetAllFacilitiesForUser(email string) ([]Facility, error) {
	keyConditionExpression := "PK = :pk AND begins_with(SK, :skPrefix)"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":pk": {
			S: aws.String("USER#" + email),
		},
		":skPrefix": {
			S: aws.String("FACILITY#"),
		},
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String("Bluebean"),
		KeyConditionExpression:    aws.String(keyConditionExpression),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := um.DB.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	facilities := make([]Facility, 0)

	for _, item := range result.Items {
		facility := Facility{
			ID:       *item["FacilityID"].S,
			Name:     *item["FacilityName"].S,
			Address:  *item["FacilityAddress"].S,
			ImageURL: *item["FacilityImageURL"].S,
		}

		facilities = append(facilities, facility)
	}

	return facilities, nil
}
