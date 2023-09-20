package main

import (
	"errors"
	"net/http"

	"bitbucket.org/nemetschek-systems/bluebean-service/internal/data"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/errorconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/messageconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/utils"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/validator"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (app *application) createFacilityHandler(c *gin.Context) {
	var input struct {
		FacilityName string `json:"name"`
		Address      string `json:"address"`
		City         string `json:"city"`
		ImageBase64  string `json:"image"`
	}

	claims, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	isAuthorized := data.AuthorizeUser(claims, data.FMRole)
	if !isAuthorized {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.InvalidJSONFormatError.Error()})
		return
	}

	facility := &data.Facility{
		Name:    input.FacilityName,
		Address: input.Address,
		City:    input.City,
	}

	v := validator.New()
	if data.ValidateFacility(v, facility); !v.Valid() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": v.Errors})

		return
	}

	imageURL, err := utils.UploadFile(input.ImageBase64, utils.FacilitiesFolder, input.FacilityName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	facility.ImageURL = imageURL

	id, err := app.models.Facilities.Insert(facility)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.FailedToInsertFacilityError.Error()})
		return
	}

	facility.ID = id.String()

	userName, nameExists := claims.(jwt.MapClaims)[utils.Name].(string)
	userEmail, emailExists := claims.(jwt.MapClaims)[utils.EmailAddress].(string)
	userRole, roleExists := claims.(jwt.MapClaims)[utils.Role].(string)
	if !nameExists || !emailExists || !roleExists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	user := &data.User{
		Name:  userName,
		Email: userEmail,
		Role:  userRole,
	}

	app.models.UserFacilities.Insert(user, facility)

	c.JSON(http.StatusCreated, facility)
}

func (app *application) getFacilityHandler(c *gin.Context) {
	facilityID := c.Param("facilityID")

	claims, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	userEmail, exists := claims.(jwt.MapClaims)[utils.EmailAddress].(string)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	facility, err := app.models.Facilities.Get(facilityID)
	if err != nil {
		switch {
		case errors.Is(err, errorconstants.RecordNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.RecordNotFoundError.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		}
		return
	}

	_, err = app.models.UserFacilities.Get(userEmail, facilityID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	c.JSON(http.StatusOK, facility)
}

type EmailData struct {
	FacilityName string
	UserRole     string
	RegisterLink string
}

func (app *application) addUserToFacilityHandler(c *gin.Context) {
	var input struct {
		FacilityID string `json:"facilityID"`
		Email      string `json:"email"`
		Role       string `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.InvalidJSONFormatError.Error()})
		return
	}

	inputRoleIsPermitted := validator.PermittedValue[string](input.Role, data.OwnerRole, data.MaintainerRole)
	if !inputRoleIsPermitted {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.RoleNotPermittedError.Error()})
		return
	}

	user, err := app.models.Users.Get(input.Email)
	if err != nil {
		//If user doesn't exist: generate a register link and send the user an email with the register link.
		facility, err := app.models.Facilities.Get(input.FacilityID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.RecordNotFoundError.Error()})
			return
		}

		registerLink := app.generateRegisterLink(input.Email, input.Role)

		emailData := EmailData{
			FacilityName: facility.Name,
			UserRole:     input.Role,
			RegisterLink: registerLink,
		}

		err = app.mailer.Send(input.Email, "user_invite.tmpl", emailData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": messageconstants.InvitationEmailSendMessage})
		return
	}

	roleIsPermitted := validator.PermittedValue[string](user.Role, data.OwnerRole, data.MaintainerRole)
	if !roleIsPermitted || input.Role != user.Role {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	//If user exists: add him to the facility and check whether he is already in the facility
	addedUser, err := app.models.Facilities.AddUserToFacility(user, input.FacilityID, app.models.Users, app.models.UserFacilities)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.UserAlreadyInFacilityError.Error()})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, addedUser)
}

func (app *application) removeUserFromFacilityHandler(c *gin.Context) {
	facilityID := c.Param("facilityID")
	email := c.Param("email")

	claims, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	isAuthorized := data.AuthorizeUser(claims, data.FMRole)
	if !isAuthorized {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	userEmail, exists := claims.(jwt.MapClaims)[utils.EmailAddress].(string)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"user": errorconstants.InvalidTokenClaimsError.Error()})
		return
	}

	_, err := app.models.UserFacilities.Get(userEmail, facilityID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"user": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	err = app.models.Facilities.RemoveUserFromFacility(email, facilityID, app.models.Users)
	if err != nil {
		switch {
		case errors.Is(err, errorconstants.RecordNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, errorconstants.UserFacilityRelashionshipError):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messageconstants.UserRemovedFromFacilityMessage})
}

func (app *application) getAllUsersForFacility(c *gin.Context) {
	facilityID := c.Param("facilityID")

	_, err := app.models.Facilities.Get(facilityID)
	if err != nil {
		switch {
		case errors.Is(err, errorconstants.RecordNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.RecordNotFoundError.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		}
		return
	}

	claims, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	isAuthorized := data.AuthorizeUser(claims, data.FMRole, data.OwnerRole)
	if !isAuthorized {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.UserIsNotAuthorizedError.Error()})

		return
	}

	users, err := app.models.Facilities.GetAllUsersForFacility(facilityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (app *application) addAssetToFacilityHandler(c *gin.Context) {
	claims, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	isAuthorized := data.AuthorizeUser(claims, data.FMRole)
	if !isAuthorized {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	var input struct {
		FacilityID string `json:"facilityID"`
		AssetName  string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.InvalidJSONFormatError.Error()})
		return
	}

	userEmail, exists := claims.(jwt.MapClaims)[utils.EmailAddress].(string)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	_, err := app.models.UserFacilities.Get(userEmail, input.FacilityID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	asset, err := app.models.Facilities.AddAssetToFacility(input.FacilityID, input.AssetName)
	if err != nil {
		switch {
		case errors.Is(err, errorconstants.RecordNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.RecordNotFoundError.Error()})
		case errors.Is(err, errorconstants.AssetAlreadyInFacilityError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.AssetAlreadyInFacilityError.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, asset)
}

func (app *application) getAllSpacesForFacility(c *gin.Context) {
	facilityID := c.Param("facilityID")

	_, err := app.models.Facilities.Get(facilityID)
	if err != nil {
		switch {
		case errors.Is(err, errorconstants.RecordNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.RecordNotFoundError.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		}
		return
	}

	claims, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	userEmail, exists := claims.(jwt.MapClaims)[utils.EmailAddress].(string)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"user": errorconstants.InvalidTokenClaimsError.Error()})
		return
	}

	_, err = app.models.UserFacilities.Get(userEmail, facilityID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"user": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	spaces, err := app.models.Facilities.GetAllSpacesForFacility(facilityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, spaces)
}

func (app *application) removeAssetFromFacilityHandler(c *gin.Context) {
	claims, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	isAuthorized := data.AuthorizeUser(claims, data.FMRole)
	if !isAuthorized {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	var input struct {
		FacilityID string `json:"facilityID"`
		AssetName  string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.InvalidJSONFormatError.Error()})
		return
	}

	userEmail, exists := claims.(jwt.MapClaims)[utils.EmailAddress].(string)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	_, err := app.models.UserFacilities.Get(userEmail, input.FacilityID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	err = app.models.Facilities.RemoveAssetFromFacility(input.FacilityID, input.AssetName)
	if err != nil {
		switch {
		case errors.Is(err, errorconstants.RecordNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.RecordNotFoundError.Error()})
		case errors.Is(err, errorconstants.AssetNotInFacilityError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.AssetNotInFacilityError.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messageconstants.AssetRemovedFromFacilityMessage})
}
