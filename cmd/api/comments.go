package main

import (
	"errors"
	"net/http"
	"time"

	"bitbucket.org/nemetschek-systems/bluebean-service/internal/data"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/errorconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/utils"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (app *application) createCommentHandler(c *gin.Context) {
	var input struct {
		PunchID    string `json:"punchID"`
		SpaceID    string `json:"spaceID"`
		FacilityID string `json:"facilityID"`
		Text       string `json:"text"`
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

	userName, exists := claims.(jwt.MapClaims)[utils.Name].(string)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	_, err := app.models.Facilities.Get(input.FacilityID)
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

	_, err = app.models.Punches.Get(input.PunchID, input.FacilityID, input.SpaceID)
	if err != nil {
		switch {
		case errors.Is(err, errorconstants.RecordNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.RecordNotFoundError.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		}
		return
	}

	comment := &data.Comment{
		PunchID:      input.PunchID,
		SpaceID:      input.SpaceID,
		FacilityID:   input.FacilityID,
		Text:         input.Text,
		CreatedOn:    time.Now().UTC().Format(time.RFC3339),
		CreatorEmail: userEmail,
		CreatorName:  userName,
	}

	v := validator.New()
	if data.ValidateComment(v, comment); !v.Valid() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": v.Errors})
		return
	}

	commentId, err := app.models.Comments.Insert(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.FailedToInsertCommentError.Error()})
		return
	}

	comment.ID = commentId.String()

	c.JSON(http.StatusCreated, comment)
}

func (app *application) getAllCommentsForPunchHandler(c *gin.Context) {
	facilityID := c.Param("facilityID")
	spaceID := c.Param("spaceID")
	punchID := c.Param("punchID")

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

	_, err = app.models.Punches.Get(punchID, facilityID, spaceID)
	if err != nil {
		switch {
		case errors.Is(err, errorconstants.RecordNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.RecordNotFoundError.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		}
		return
	}

	comments, err := app.models.Comments.GetAllCommentsForPunch(punchID, spaceID, facilityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}
