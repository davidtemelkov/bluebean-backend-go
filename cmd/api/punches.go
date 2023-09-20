package main

import (
	"errors"
	"net/http"
	"regexp"

	"bitbucket.org/nemetschek-systems/bluebean-service/internal/data"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/errorconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/generalconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/messageconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/utils"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (app *application) createPunchHandler(c *gin.Context) {
	var input struct {
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
		Creator     string `json:"creator"`
		Asset       string `json:"asset"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.InvalidJSONFormatError.Error()})
		return
	}

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

	facility, err := app.models.Facilities.Get(input.FacilityID)
	if err != nil {
		switch {
		case errors.Is(err, errorconstants.RecordNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.RecordNotFoundError.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		}
		return
	}

	_, err = app.models.UserFacilities.Get(userEmail, input.FacilityID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	_, err = app.models.Spaces.Get(input.SpaceID, input.FacilityID)
	if err != nil {
		switch {
		case errors.Is(err, errorconstants.RecordNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.RecordNotFoundError.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		}
		return
	}

	punch := &data.Punch{
		FacilityID:  input.FacilityID,
		SpaceID:     input.SpaceID,
		Title:       input.Title,
		Description: input.Description,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		CoordX:      input.CoordX,
		CoordY:      input.CoordY,
		Status:      input.Status,
		Creator:     userEmail,
		Asset:       input.Asset,
	}

	punch.Assignee = input.Assignee
	punch.Status = generalconstants.StatusInProgress
	if punch.Assignee == "" {
		punch.Assignee = generalconstants.StatusUnassigned
		punch.Status = generalconstants.StatusUnassigned
	}

	punch.Asset = input.Asset
	if punch.Asset == "" {
		punch.Asset = generalconstants.AssetNone
	}

	v := validator.New()
	if data.ValidatePunch(v, punch); !v.Valid() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": v.Errors})
		return
	}

	dateTimePattern := regexp.MustCompile(generalconstants.ISO8601)
	if !validator.Matches(input.StartDate, dateTimePattern) || !validator.Matches(input.EndDate, dateTimePattern) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InvalidDateTimeFormatError.Error()})
		return
	}

	if !v.IsValidDateTimeRange(input.StartDate, input.EndDate) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InvalidDateTimeRangeError.Error()})
		return
	}

	permittedStatuses := []string{generalconstants.StatusUnassigned, generalconstants.StatusInProgress, generalconstants.StatusCompleted}
	if !validator.PermittedValue(input.Status, permittedStatuses...) {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.InvalidPunchStatusError.Error()})
		return
	}

	maintainers := facility.Maintainers
	if punch.Assignee != generalconstants.StatusUnassigned && !validator.PermittedValue(punch.Assignee, maintainers...) {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.AssigneeIsNotMaintainerError.Error()})
		return
	}

	if _, exists := facility.Assets[punch.Asset]; !exists && punch.Asset != generalconstants.AssetNone {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.AssetNotInFacilityError.Error()})
		return
	}

	punch.Creator = userEmail

	punchId, err := app.models.Punches.Insert(punch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.FailedToInsertPunchError.Error()})
		return
	}

	punch.ID = punchId.String()

	c.JSON(http.StatusCreated, punch)
}

func (app *application) getPunchHandler(c *gin.Context) {
	punchID := c.Param("punchID")
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

	punch, err := app.models.Punches.Get(punchID, facilityID, spaceID)
	if err != nil {
		switch {
		case errors.Is(err, errorconstants.RecordNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.RecordNotFoundError.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, punch)
}

func (app *application) getAllPunchesForSpaceHandler(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	_, err = app.models.UserFacilities.Get(userEmail, facilityID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"user": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	_, err = app.models.Spaces.Get(spaceID, facilityID)
	if err != nil {
		switch {
		case errors.Is(err, errorconstants.RecordNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.RecordNotFoundError.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		}
		return
	}

	punches, err := app.models.Punches.GetAllPunchesForSpace(spaceID, facilityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, punches)
}

func (app *application) getAllPunchesForFacilityHandler(c *gin.Context) {
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

	userEmail, userEmailExists := claims.(jwt.MapClaims)[utils.EmailAddress].(string)
	if !userEmailExists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	_, err = app.models.UserFacilities.Get(userEmail, facilityID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"user": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	punches, err := app.models.Punches.GetAllPunchesForFacility(facilityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}
	c.JSON(http.StatusOK, punches)
}

func (app *application) editPunchHandler(c *gin.Context) {
	var input struct {
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
		Creator     string `json:"creator"`
		Asset       string `json:"asset"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.InvalidJSONFormatError.Error()})
		return
	}

	_, err := app.models.Punches.Get(input.ID, input.FacilityID, input.SpaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.PunchNotExistError.Error()})
		return
	}

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

	facility, err := app.models.Facilities.Get(input.FacilityID)
	if err != nil {
		switch {
		case errors.Is(err, errorconstants.RecordNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.RecordNotFoundError.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		}
		return
	}

	_, err = app.models.UserFacilities.Get(userEmail, input.FacilityID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"user": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	_, err = app.models.Spaces.Get(input.SpaceID, input.FacilityID)
	if err != nil {
		switch {
		case errors.Is(err, errorconstants.RecordNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.RecordNotFoundError.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		}
		return
	}

	punch := &data.Punch{
		ID:          input.ID,
		FacilityID:  input.FacilityID,
		SpaceID:     input.SpaceID,
		Creator:     input.Creator,
		Title:       input.Title,
		Description: input.Description,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		CoordX:      input.CoordX,
		CoordY:      input.CoordY,
		Status:      input.Status,
		Asset:       input.Asset,
	}

	punch.Assignee = input.Assignee
	if punch.Assignee == "" {
		punch.Assignee = generalconstants.StatusUnassigned
	}

	punch.Asset = input.Asset
	if punch.Asset == "" {
		punch.Asset = generalconstants.AssetNone
	}

	v := validator.New()
	if data.ValidatePunch(v, punch); !v.Valid() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": v.Errors})
		return
	}

	dateTimePattern := regexp.MustCompile(generalconstants.ISO8601)
	if !validator.Matches(input.StartDate, dateTimePattern) || !validator.Matches(input.EndDate, dateTimePattern) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InvalidDateTimeFormatError.Error()})
		return
	}

	if !v.IsValidDateTimeRange(input.StartDate, input.EndDate) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InvalidDateTimeRangeError.Error()})
		return
	}

	permittedStatuses := []string{generalconstants.StatusUnassigned, generalconstants.StatusInProgress, generalconstants.StatusCompleted}
	if !validator.PermittedValue(input.Status, permittedStatuses...) {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.InvalidPunchStatusError.Error()})
		return
	}

	maintainers := facility.Maintainers
	if punch.Assignee != generalconstants.StatusUnassigned && !validator.PermittedValue(punch.Assignee, maintainers...) {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.AssigneeIsNotMaintainerError.Error()})
		return
	}

	if _, exists := facility.Assets[punch.Asset]; !exists && punch.Asset != generalconstants.AssetNone {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.AssetNotInFacilityError.Error()})
		return
	}

	err = app.models.Punches.Edit(punch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	c.JSON(http.StatusCreated, punch)
}

func (app *application) deletePunchHandler(c *gin.Context) {
	punchID := c.Param("punchID")
	facilityID := c.Param("facilityID")
	spaceID := c.Param("spaceID")

	punch, err := app.models.Punches.Get(punchID, facilityID, spaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.PunchNotExistError.Error()})
		return
	}

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

	userRole, exists := claims.(jwt.MapClaims)[utils.EmailAddress].(string)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	// User must be the FM of the facility or the creator of the punch
	_, err = app.models.UserFacilities.Get(userEmail, facilityID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"user": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	roleIsPermitted := validator.PermittedValue[string](userRole, data.FMRole)
	if !roleIsPermitted && userEmail != punch.Creator {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	err = app.models.Punches.Delete(punchID, facilityID, spaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messageconstants.PunchDeletedSuccessfullyMessage})
}
