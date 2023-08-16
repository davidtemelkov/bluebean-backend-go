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
		ImageUrl     string   `json:"imageurl"`
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
		ImageUrl: input.ImageUrl,
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
	id := c.Param("id")
	facility, err := app.models.Facilities.Get(id)
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