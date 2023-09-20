package main

import (
	"errors"
	"net/http"

	"bitbucket.org/nemetschek-systems/bluebean-service/internal/data"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/errorconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/utils"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (app *application) createSpaceHandler(c *gin.Context) {
	var input struct {
		FacilityID   string `json:"facilityId"`
		Name         string `json:"name"`
		Location     string `json:"location"`
		SchemaBase64 string `json:"schema"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
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

	space := &data.Space{
		FacilityID: input.FacilityID,
		Name:       input.Name,
		Location:   input.Location,
	}

	v := validator.New()
	if data.ValidateSpace(v, space); !v.Valid() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": v.Errors})
		return
	}

	schemaURL, err := utils.UploadFile(input.SchemaBase64, utils.SpacesFolder, input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	space.SchemaURL = schemaURL

	id, err := app.models.Spaces.Insert(space)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.FailedToInsertSpaceError.Error()})
		return
	}

	space.ID = id.String()

	c.JSON(http.StatusCreated, space)
}

func (app *application) getSpaceHandler(c *gin.Context) {
	facilityID := c.Param("facilityID")
	spaceID := c.Param("spaceID")

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

	space, err := app.models.Spaces.Get(spaceID, facilityID)
	if err != nil {
		switch {
		case errors.Is(err, errorconstants.RecordNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.RecordNotFoundError.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, space)
}
