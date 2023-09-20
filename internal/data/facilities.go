package data

import (
	"context"
	"time"

	"bitbucket.org/nemetschek-systems/bluebean-service/internal/errorconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/generalconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/validator"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
)

type Facility struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Address     string            `json:"address"`
	City        string            `json:"city"`
	Owners      []string          `json:"owners"`
	Maintainers []string          `json:"maintainers"`
	Assets      map[string]string `json:"assets"`
	ImageURL    string            `json:"imageURL"`
}
type FacilityModel struct {
	DB *dynamodb.DynamoDB
}

func ValidateFacility(v *validator.Validator, facility *Facility) {
	v.Check(facility.Name != "", "name", errorconstants.RequiredFieldError.Error())
	v.Check(len(facility.Name) >= 2, "name", errorconstants.NameMinLengthError.Error())
	v.Check(len(facility.Name) < 50, "name", errorconstants.NameMaxLengthError.Error())
	v.Check(facility.Address != "", "address", errorconstants.RequiredFieldError.Error())
	v.Check(len(facility.Address) > 5, "address", errorconstants.AddressMinLengthError.Error())
	v.Check(len(facility.Address) < 100, "address", errorconstants.AddressMaxLengthError.Error())
	v.Check(facility.City != "", "city", errorconstants.RequiredFieldError.Error())
	v.Check(len(facility.City) > 2, "city", errorconstants.CityMinLengthError.Error())
	v.Check(len(facility.City) < 50, "city", errorconstants.CityMaxLengthError.Error())
}

func (fm FacilityModel) Insert(facility *Facility) (uuid.UUID, error) {
	id := uuid.New()

	item := map[string]*dynamodb.AttributeValue{
		generalconstants.PK: {
			S: aws.String(
				generalconstants.FacilityPrefix + id.String(),
			),
		},
		generalconstants.SK: {
			S: aws.String(
				generalconstants.FacilityPrefix + id.String(),
			),
		},
		"ID": {
			S: aws.String(id.String()),
		},
		"Name": {
			S: aws.String(facility.Name),
		},
		"Address": {
			S: aws.String(facility.Address),
		},
		"City": {
			S: aws.String(facility.City),
		},
		"ImageURL": {
			S: aws.String(facility.ImageURL),
		},
		"Assets": {
			M: make(map[string]*dynamodb.AttributeValue),
		},
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(generalconstants.TableName),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := fm.DB.PutItemWithContext(ctx, input)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (fm FacilityModel) Get(id string) (*Facility, error) {
	if id == "" {
		return nil, errorconstants.RecordNotFoundError
	}

	keyCondition := expression.
		Key(generalconstants.PK).
		Equal(expression.Value(generalconstants.FacilityPrefix + id)).
		And(expression.Key(generalconstants.SK).Equal(expression.Value(generalconstants.FacilityPrefix + id)))

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

	result, err := fm.DB.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, errorconstants.RecordNotFoundError
	}

	facility := &Facility{}
	err = dynamodbattribute.UnmarshalMap(result.Items[0], facility)
	if err != nil {
		return nil, err
	}

	return facility, nil
}

type AddedUser struct {
	FacilityID  string `json:"facilityId"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	UserAddedOn string `json:"userAddedOn"`
}

func (fm FacilityModel) AddUserToFacility(user *User, facilityID string, um UserModel, ufm UserFacilityModel) (*AddedUser, error) {
	facility, err := fm.Get(facilityID)
	if err != nil {
		return nil, errorconstants.RecordNotFoundError
	}

	err = ufm.Insert(user, facility)
	if err != nil {
		return nil, errorconstants.InternalServerError
	}

	err = fm.AddUserToFacilityRoleSet(user.Email, user.Role, facilityID)
	if err != nil {
		return nil, errorconstants.InternalServerError
	}

	addedUser := &AddedUser{
		FacilityID:  facilityID,
		Name:        user.Name,
		Email:       user.Email,
		Role:        user.Role,
		UserAddedOn: time.Now().UTC().Format(time.RFC3339),
	}

	return addedUser, nil
}

func (fm FacilityModel) AddUserToFacilityRoleSet(userEmail, role, facilityID string) error {
	var updateExpression string
	switch role {
	case OwnerRole:
		updateExpression = "ADD Owners :userEmail"
	case MaintainerRole:
		updateExpression = "ADD Maintainers :userEmail"
	default:
		return errorconstants.RoleNotPermittedError
	}

	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":userEmail": {
			SS: []*string{aws.String(userEmail)},
		},
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(generalconstants.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			generalconstants.PK: {
				S: aws.String(
					generalconstants.FacilityPrefix + facilityID,
				),
			},
			generalconstants.SK: {
				S: aws.String(
					generalconstants.FacilityPrefix + facilityID,
				),
			},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := fm.DB.UpdateItemWithContext(ctx, updateInput)
	if err != nil {
		return err
	}

	return nil
}

func (fm FacilityModel) RemoveUserFromFacility(userEmail, facilityID string, um UserModel) error {
	user, err := um.Get(userEmail)
	if err != nil {
		return errorconstants.RecordNotFoundError
	}

	_, err = fm.Get(facilityID)
	if err != nil {
		return errorconstants.RecordNotFoundError
	}

	userFacilityKey := map[string]*dynamodb.AttributeValue{
		generalconstants.PK: {
			S: aws.String(
				generalconstants.UserPrefix + userEmail,
			),
		},
		generalconstants.SK: {
			S: aws.String(
				generalconstants.FacilityPrefix + facilityID,
			),
		},
	}

	result, err := fm.DB.GetItemWithContext(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(generalconstants.TableName),
		Key:       userFacilityKey,
	})
	if err != nil || result.Item == nil {
		return errorconstants.UserFacilityRelashionshipError
	}

	_, err = fm.DB.DeleteItemWithContext(context.Background(), &dynamodb.DeleteItemInput{
		TableName: aws.String(generalconstants.TableName),
		Key:       userFacilityKey,
	})
	if err != nil {
		return err
	}

	roleIsPermitted := validator.PermittedValue[string](user.Role, OwnerRole, MaintainerRole)
	if !roleIsPermitted {
		return errorconstants.RoleNotPermittedError
	}

	err = fm.RemoveUserFromFacilityRoleSet(userEmail, user.Role, facilityID)
	if err != nil {
		return errorconstants.InternalServerError
	}

	return nil
}

func (fm FacilityModel) RemoveUserFromFacilityRoleSet(userEmail, userRole, id string) error {
	var updateExpression string
	switch userRole {
	case OwnerRole:
		updateExpression = "DELETE Owners :userEmail"
	case MaintainerRole:
		updateExpression = "DELETE Maintainers :userEmail"
	default:
		return errorconstants.RoleNotPermittedError
	}

	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":userEmail": {SS: []*string{aws.String(userEmail)}},
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(generalconstants.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			generalconstants.PK: {
				S: aws.String(
					generalconstants.FacilityPrefix + id,
				),
			},
			generalconstants.SK: {
				S: aws.String(
					generalconstants.FacilityPrefix + id,
				),
			},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := fm.DB.UpdateItemWithContext(ctx, updateInput)
	if err != nil {
		return err
	}

	return nil
}

func (fm FacilityModel) GetAllUsersForFacility(id string) ([]User, error) {
	if id == "" {
		return nil, errorconstants.RecordNotFoundError
	}

	gsi1pk := generalconstants.FacilityPrefix + id
	gsi1skPrefix := generalconstants.UserPrefix

	keyCondition := expression.Key(generalconstants.GSI1PK).Equal(expression.Value(gsi1pk)).
		And(expression.Key(generalconstants.GSI1SK).BeginsWith(gsi1skPrefix))

	builder, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		return nil, err
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(generalconstants.TableName),
		IndexName:                 aws.String(generalconstants.GSI1),
		KeyConditionExpression:    builder.KeyCondition(),
		ExpressionAttributeNames:  builder.Names(),
		ExpressionAttributeValues: builder.Values(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := fm.DB.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)

	for _, item := range result.Items {
		user := User{
			Email:   *item["UserEmail"].S,
			Name:    *item["UserName"].S,
			Role:    *item["UserRole"].S,
			AddedOn: *item["UserAddedOn"].S,
		}
		users = append(users, user)
	}

	return users, nil
}

func (fm FacilityModel) GetAllSpacesForFacility(id string) ([]Space, error) {
	if id == "" {
		return nil, errorconstants.RecordNotFoundError
	}

	pk := generalconstants.FacilityPrefix + id
	skPrefix := generalconstants.SpacePrefix

	keyCondition := expression.Key("PK").Equal(expression.Value(pk)).
		And(expression.Key("SK").BeginsWith(skPrefix))

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

	result, err := fm.DB.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	spaces := make([]Space, 0)

	for _, item := range result.Items {
		space := Space{
			ID:         *item["ID"].S,
			Name:       *item["Name"].S,
			Location:   *item["Location"].S,
			SchemaURL:  *item["SchemaURL"].S,
			FacilityID: *item["FacilityID"].S,
		}
		spaces = append(spaces, space)
	}

	return spaces, nil
}

type Asset struct {
	Name    string `json:"name"`
	AddedOn string `json:"addedOn"`
}

func (fm FacilityModel) AddAssetToFacility(facilityID, assetName string) (*Asset, error) {
	facility, err := fm.Get(facilityID)
	if err != nil {
		return nil, errorconstants.RecordNotFoundError
	}

	timeNow := time.Now().UTC().Format(time.RFC3339)

	updateExpression := "SET Assets.#key = :assetValue"
	conditionExpression := "attribute_not_exists(Assets.#key)"
	expressionAttributeNames := map[string]*string{
		"#key": aws.String(assetName),
	}
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":assetValue": {
			S: aws.String(timeNow),
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(generalconstants.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			generalconstants.PK: {
				S: aws.String(
					generalconstants.FacilityPrefix + facilityID,
				),
			},
			generalconstants.SK: {
				S: aws.String(
					generalconstants.FacilityPrefix + facilityID,
				),
			},
		},
		UpdateExpression:          aws.String(updateExpression),
		ConditionExpression:       aws.String(conditionExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = fm.DB.UpdateItemWithContext(ctx, input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			return nil, errorconstants.AssetAlreadyInFacilityError
		}
		return nil, err
	}

	if facility.Assets == nil {
		facility.Assets = make(map[string]string)
	}

	facility.Assets[assetName] = timeNow

	asset := &Asset{Name: assetName, AddedOn: timeNow}

	return asset, nil
}

func (fm FacilityModel) RemoveAssetFromFacility(facilityID, assetName string) error {
	facility, err := fm.Get(facilityID)
	if err != nil {
		return errorconstants.RecordNotFoundError
	}

	if _, exists := facility.Assets[assetName]; !exists {
		return errorconstants.AssetNotInFacilityError
	}

	updateExpression := "REMOVE Assets.#key"
	expressionAttributeNames := map[string]*string{
		"#key": aws.String(assetName),
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(generalconstants.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			generalconstants.PK: {
				S: aws.String(
					generalconstants.FacilityPrefix + facilityID,
				),
			},
			generalconstants.SK: {
				S: aws.String(
					generalconstants.FacilityPrefix + facilityID,
				),
			},
		},
		UpdateExpression:         aws.String(updateExpression),
		ExpressionAttributeNames: expressionAttributeNames,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = fm.DB.UpdateItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	delete(facility.Assets, assetName)

	return nil
}
