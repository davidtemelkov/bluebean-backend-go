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

type Punch struct {
	ID          string `json:"id"`
	FacilityID  string `json:"facilityID"`
	SpaceID     string `json:"spaceID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	CoordX      string `json:"coordX"`
	CoordY      string `json:"coordY"`
	Status      string `json:"status"`
	Assignee    string `json:"assignee"`
	Creator     string `json:"creator,omitempty"`
	Asset       string `json:"asset"`
	GSI1PK      string `json:"GSI1PK,omitempty"`
	GSI1SK      string `json:"GSI1SK,omitempty"`
}

type PunchModel struct {
	DB *dynamodb.DynamoDB
}

var (
	UnassignedStatus = "Unassigned"
	InProgressStatus = "In progress"
	CompletedStatus  = "Completed"
)

func ValidatePunch(v *validator.Validator, punch *Punch) {
	v.Check(punch.Title != "", "title", errorconstants.RequiredFieldError.Error())
	v.Check(len(punch.Title) >= 5, "title", errorconstants.PunchTitleMinLengthError.Error())
	v.Check(len(punch.Title) < 100, "title", errorconstants.PunchTitleMaxLengthError.Error())
	v.Check(len(punch.Description) < 500, "description", errorconstants.PunchDescriptionMaxLengthError.Error())
	v.Check(punch.StartDate != "", "startDate", errorconstants.RequiredFieldError.Error())
	v.Check(punch.EndDate != "", "endDate", errorconstants.RequiredFieldError.Error())
	v.Check(punch.CoordX != "", "coordX", errorconstants.RequiredFieldError.Error())
	v.Check(len(punch.CoordX) >= 0, "coordX", errorconstants.PunchCoordXMinValueError.Error())
	v.Check(len(punch.CoordX) <= 100, "coordX", errorconstants.PunchCoordXMaxValueError.Error())
	v.Check(punch.CoordY != "", "coordY", errorconstants.RequiredFieldError.Error())
	v.Check(len(punch.CoordY) >= 0, "coordY", errorconstants.PunchCoordYMinValueError.Error())
	v.Check(len(punch.CoordY) <= 100, "coordY", errorconstants.PunchCoordYMaxValueError.Error())
	v.Check(punch.Status != "", "status", errorconstants.RequiredFieldError.Error())
}

func (pm PunchModel) Insert(punch *Punch) (uuid.UUID, error) {
	id := uuid.New()

	item := map[string]*dynamodb.AttributeValue{
		generalconstants.PK: {
			S: aws.String(
				generalconstants.FacilityPrefix + punch.FacilityID +
					generalconstants.SpacePrefix + punch.SpaceID),
		},
		generalconstants.SK: {
			S: aws.String(
				generalconstants.PunchSKPrefix + id.String(),
			),
		},
		"ID": {
			S: aws.String(id.String()),
		},
		"FacilityID": {
			S: aws.String(punch.FacilityID),
		},
		"SpaceID": {
			S: aws.String(punch.SpaceID),
		},
		"Title": {
			S: aws.String(punch.Title),
		},
		"Description": {
			S: aws.String(punch.Description),
		},
		"StartDate": {
			S: aws.String(punch.StartDate),
		},
		"EndDate": {
			S: aws.String(punch.EndDate),
		},
		"CoordX": {
			S: aws.String(punch.CoordX),
		},
		"CoordY": {
			S: aws.String(punch.CoordY),
		},
		"Status": {
			S: aws.String(punch.Status),
		},
		"Assignee": {
			S: aws.String(punch.Assignee),
		},
		"Creator": {
			S: aws.String(punch.Creator),
		},
		"Asset": {
			S: aws.String(punch.Asset),
		},
		"GSI1PK": {
			S: aws.String(
				generalconstants.FacilityPrefix + punch.FacilityID,
			),
		},
		"GSI1SK": {
			S: aws.String(
				generalconstants.PunchSKPrefix + id.String(),
			),
		},
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(generalconstants.TableName),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := pm.DB.PutItemWithContext(ctx, input)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (pm PunchModel) Get(punchID, facilityID, spaceID string) (*Punch, error) {
	if punchID == "" || facilityID == "" || spaceID == "" {
		return nil, errorconstants.RecordNotFoundError
	}

	pk := generalconstants.FacilityPrefix + facilityID + generalconstants.SpacePrefix + spaceID
	sk := generalconstants.PunchSKPrefix + punchID

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

	result, err := pm.DB.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, errorconstants.RecordNotFoundError
	}

	item := result.Items[0]
	punch := &Punch{
		ID:          *item["ID"].S,
		FacilityID:  *item["FacilityID"].S,
		SpaceID:     *item["SpaceID"].S,
		Title:       *item["Title"].S,
		Description: *item["Description"].S,
		StartDate:   *item["StartDate"].S,
		EndDate:     *item["EndDate"].S,
		CoordX:      *item["CoordX"].S,
		CoordY:      *item["CoordY"].S,
		Status:      *item["Status"].S,
		Assignee:    *item["Assignee"].S,
		Creator:     *item["Creator"].S,
		Asset:       *item["Asset"].S,
	}

	return punch, nil
}

func (pm PunchModel) GetAllPunchesForSpace(spaceID, facilityID string) ([]Punch, error) {
	if spaceID == "" || facilityID == "" {
		return nil, errorconstants.RecordNotFoundError
	}

	pk := generalconstants.FacilityPrefix + facilityID + generalconstants.SpacePrefix + spaceID
	skPrefix := generalconstants.PunchSKPrefix

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

	result, err := pm.DB.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	punches := make([]Punch, 0)

	for _, item := range result.Items {
		punch := Punch{
			ID:          *item["ID"].S,
			FacilityID:  *item["FacilityID"].S,
			SpaceID:     *item["SpaceID"].S,
			Title:       *item["Title"].S,
			Description: *item["Description"].S,
			StartDate:   *item["StartDate"].S,
			EndDate:     *item["EndDate"].S,
			CoordX:      *item["CoordX"].S,
			CoordY:      *item["CoordY"].S,
			Status:      *item["Status"].S,
			Assignee:    *item["Assignee"].S,
			Creator:     *item["Creator"].S,
			Asset:       *item["Asset"].S,
		}
		punches = append(punches, punch)
	}

	return punches, nil
}

func (pm PunchModel) GetAllPunchesForFacility(facilityID string) ([]Punch, error) {
	if facilityID == "" {
		return nil, errorconstants.RecordNotFoundError
	}

	gsi1pk := generalconstants.FacilityPrefix + facilityID
	gsi1skPrefix := generalconstants.PunchSKPrefix

	keyCondition := expression.Key(generalconstants.GSI1PK).Equal(expression.Value(gsi1pk)).
		And(expression.Key(generalconstants.GSI1SK).BeginsWith(gsi1skPrefix))

	builder, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		return nil, err
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(generalconstants.TableName),
		IndexName:                 aws.String("GSI1"),
		KeyConditionExpression:    builder.KeyCondition(),
		ExpressionAttributeNames:  builder.Names(),
		ExpressionAttributeValues: builder.Values(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := pm.DB.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}
	punches := make([]Punch, 0)

	if len(result.Items) > 0 {
		for _, item := range result.Items {
			punch := Punch{
				ID:          *item["ID"].S,
				FacilityID:  *item["FacilityID"].S,
				SpaceID:     *item["SpaceID"].S,
				Title:       *item["Title"].S,
				Description: *item["Description"].S,
				StartDate:   *item["StartDate"].S,
				EndDate:     *item["EndDate"].S,
				CoordX:      *item["CoordX"].S,
				CoordY:      *item["CoordY"].S,
				Status:      *item["Status"].S,
				Assignee:    *item["Assignee"].S,
				Creator:     *item["Creator"].S,
				Asset:       *item["Asset"].S,
			}
			punches = append(punches, punch)
		}
	}

	return punches, nil
}

func (pm PunchModel) Edit(updatedPunch *Punch) error {
	builder := expression.NewBuilder()

	updateExpression := expression.Set(
		expression.Name("Title"),
		expression.Value(updatedPunch.Title),
	).Set(
		expression.Name("Description"),
		expression.Value(updatedPunch.Description),
	).Set(
		expression.Name("StartDate"),
		expression.Value(updatedPunch.StartDate),
	).Set(
		expression.Name("EndDate"),
		expression.Value(updatedPunch.EndDate),
	).Set(
		expression.Name("CoordX"),
		expression.Value(updatedPunch.CoordX),
	).Set(
		expression.Name("CoordY"),
		expression.Value(updatedPunch.CoordY),
	).Set(
		expression.Name("Status"),
		expression.Value(updatedPunch.Status),
	).Set(
		expression.Name("Assignee"),
		expression.Value(updatedPunch.Assignee),
	).Set(
		expression.Name("Asset"),
		expression.Value(updatedPunch.Asset))

	builder = builder.WithUpdate(updateExpression)

	expr, err := builder.Build()
	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(generalconstants.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			generalconstants.PK: {
				S: aws.String(
					generalconstants.FacilityPrefix + updatedPunch.FacilityID +
						generalconstants.SpacePrefix + updatedPunch.SpaceID,
				),
			},
			generalconstants.SK: {
				S: aws.String(
					generalconstants.PunchSKPrefix + updatedPunch.ID,
				),
			},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              aws.String("ALL_NEW"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = pm.DB.UpdateItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (pm PunchModel) Delete(punchID, facilityID, spaceID string) error {
	writeRequests := make([]*dynamodb.WriteRequest, 0)

	skPrefix := generalconstants.PunchSKPrefix + punchID

	keyCondition := expression.Key(generalconstants.PK).
		Equal(expression.Value(
			generalconstants.FacilityPrefix + facilityID +
				generalconstants.SpacePrefix + spaceID,
		)).
		And(expression.Key(generalconstants.SK).
			BeginsWith(skPrefix))

	builder, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		return err
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(generalconstants.TableName),
		KeyConditionExpression:    builder.KeyCondition(),
		ExpressionAttributeNames:  builder.Names(),
		ExpressionAttributeValues: builder.Values(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := pm.DB.QueryWithContext(ctx, queryInput)
	if err != nil {
		return err
	}

	for _, item := range result.Items {
		writeRequests = append(writeRequests, &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]*dynamodb.AttributeValue{
					generalconstants.PK: item[generalconstants.PK],
					generalconstants.SK: item[generalconstants.SK],
				},
			},
		})
	}

	batchInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			generalconstants.TableName: writeRequests,
		},
	}

	_, err = pm.DB.BatchWriteItemWithContext(ctx, batchInput)
	if err != nil {
		return err
	}

	return nil
}
