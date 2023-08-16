package data

import (
	"errors"
	"example/bluebean-go/internal/validator"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password password `json:"-"`
	Role     string   `json:"role"`
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
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")
	// Call the standalone ValidateEmail() helper.
	ValidateEmail(v, user.Email)
	// If the plaintext password is not nil, call the standalone
	// ValidatePasswordPlaintext() helper.
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user). It's a
	// useful sanity check to include here, but it's not a problem with the data
	// provided by the client. So rather than adding an error to the validation map we
	// raise a panic instead.
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

type UserModel struct {
	DB *dynamodb.DynamoDB
}

func (m UserModel) Insert(user *User) error {
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
		ConditionExpression: aws.String("attribute_not_exists(email)"), // Prevent duplicate email insertion
	}

	_, err := m.DB.PutItem(input)
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
