package main

import (
	"errors"
	"example/bluebean-go/internal/data"
	"example/bluebean-go/internal/validator"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) createFacilityHandler(c *gin.Context) {
	var input struct {
		FacilityName string   `json:"name"`
		Owners       []string `json:"owners"`
		Address      string   `json:"address"`
		City         string   `json:"city"`
		Creator      string   `json:"creator"`
		ImageURL     string   `json:"imageurl"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	facility := &data.Facility{
		Name:     input.FacilityName,
		Owners:   input.Owners,
		Address:  input.Address,
		City:     input.City,
		Creator:  input.Creator,
		ImageURL: input.ImageURL,
	}

	v := validator.New()
	if data.ValidateFacility(v, facility); !v.Valid() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": v.Errors})
		return
	}

	if err := app.models.Facilities.InsertFacility(facility); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert facility"})
		return
	}

	c.Header("Location", fmt.Sprintf("/facilities/%s", facility.Name))
	c.JSON(http.StatusCreated, gin.H{"facility": facility})
}

func (app *application) getFacilityHandler(c *gin.Context) {
	facilityName := c.Param("facilityName")
	facility, err := app.models.Facilities.Get(facilityName)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Facility not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"facility": facility})
}

func (app *application) addUserToFacilityHandler(c *gin.Context) {
	facilityName := c.Param("facilityName")
	email := c.Param("email")

	err := app.models.Facilities.AddUserToFacility(email, facilityName, app.models.Users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User added to facility"})
}

func (app *application) getAllUsersForFacility(c *gin.Context) {
	facilityName := c.Param("facilityName")
	users, err := app.models.Facilities.GetAllUsersForFacility(facilityName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
