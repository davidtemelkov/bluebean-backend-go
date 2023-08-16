package main

import (
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
