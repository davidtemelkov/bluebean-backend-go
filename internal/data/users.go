package data

import (
	"context"
	"errors"
	"strings"
	"time"

	"bitbucket.org/nemetschek-systems/bluebean-service/internal/errorconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/generalconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/utils"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/validator"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password Password `json:"-"`
	Role     string   `json:"role"`
	AddedOn  string   `json:"addedOn,omitempty"`
}

var (
	MaintainerRole = "Maintainer"
	OwnerRole      = "Owner"
	FMRole         = "FM"
)

type Password struct {
	plaintext *string
	hash      []byte
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
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
	v.Check(email != "", "email", errorconstants.RequiredFieldError.Error())
	v.Check(validator.Matches(email, validator.EmailRX), "email", errorconstants.EmailFormatError.Error())
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", errorconstants.RequiredFieldError.Error())
	v.Check(len(password) >= 8, "password", errorconstants.PasswordMinLengthError.Error())
	v.Check(len(password) <= 72, "password", errorconstants.PasswordMaxLengthError.Error())
}

func ValidateRegisterInput(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", errorconstants.RequiredFieldError.Error())
	v.Check(len(user.Name) >= 5, "name", errorconstants.UserNameMinLengthError.Error())
	v.Check(len(user.Name) <= 50, "name", errorconstants.UserNameMaxLengthError.Error())
	v.Check(len(strings.Split(user.Name, " ")) == 2, "name", errorconstants.UserNameNoWhitespaceError.Error())

	ValidateEmail(v, user.Email)
	ValidatePasswordPlaintext(v, *user.Password.plaintext)

	v.Check(user.Role != "", "role", errorconstants.RequiredFieldError.Error())
	roleIsPermitted := validator.PermittedValue[string](user.Role, OwnerRole, MaintainerRole)
	if !roleIsPermitted {
		v.AddError("role", errorconstants.RoleNotPermittedError.Error())
	}
}

func ValidateLoginInput(v *validator.Validator, email, password string) {
	ValidateEmail(v, email)
	ValidatePasswordPlaintext(v, password)
}

type UserModel struct {
	DB *dynamodb.DynamoDB
}

func (um UserModel) Insert(user *User) error {
	item := map[string]*dynamodb.AttributeValue{
		generalconstants.PK: {
			S: aws.String(
				generalconstants.UserPrefix + user.Email,
			),
		},
		generalconstants.SK: {
			S: aws.String(
				generalconstants.UserPrefix + user.Email,
			),
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
		TableName:           aws.String(generalconstants.TableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := um.DB.PutItemWithContext(ctx, input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
				return errorconstants.DuplicateEmailError
			}
		}
		return err
	}

	return nil
}

func (um UserModel) Get(email string) (*User, error) {
	if email == "" {
		return nil, errorconstants.UserNotFoundError
	}

	pk := generalconstants.UserPrefix + email
	sk := generalconstants.UserPrefix + email

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

	result, err := um.DB.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, errorconstants.UserNotFoundError
	}

	item := result.Items[0]
	user := &User{
		Name:  *item["Name"].S,
		Email: *item["Email"].S,
		Role:  *item["Role"].S,
		Password: Password{
			hash: []byte(*item["HashedPassword"].S),
		},
	}

	return user, nil
}

func (um UserModel) CanLoginUser(password string, user *User) (bool, error) {
	passwordIsCorrect, err := user.Password.Matches(password)
	if err != nil || !passwordIsCorrect {
		return false, err
	}

	return true, nil
}

func (um UserModel) GetAllFacilitiesForUser(email string) ([]Facility, error) {
	if email == "" {
		return nil, errorconstants.UserNotFoundError
	}

	pk := generalconstants.UserPrefix + email
	skPrefix := generalconstants.FacilityPrefix

	keyCondition := expression.Key(generalconstants.PK).Equal(expression.Value(pk)).
		And(expression.Key(generalconstants.SK).BeginsWith(skPrefix))

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
			City:     *item["FacilityCity"].S,
			ImageURL: *item["FacilityImageURL"].S,
		}

		facilities = append(facilities, facility)
	}

	return facilities, nil
}

func AuthorizeUser(claims any, premittedRoles ...string) bool {
	userRole, exists := claims.(jwt.MapClaims)[utils.Role].(string)
	if !exists {
		return false
	}

	roleIsPermitted := validator.PermittedValue[string](userRole, premittedRoles...)
	if !roleIsPermitted {
		return false
	}

	return true
}
