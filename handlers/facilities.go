package handlers

import (
	"net/http"
	//"errors"
	"example/bluebean-go/database"
	exportdtos "example/bluebean-go/dtos/export"
	"example/bluebean-go/model"

	"github.com/gin-gonic/gin"
	//"gorm.io/gorm"
)

func GetOneFacility(c *gin.Context) {
	id := c.Query("facilityId")

	var facility model.Facility
	result := database.Db.First(&facility, id)
	if result.Error != nil {
		panic(result.Error)
	}

	var city model.City
	result = database.Db.First(&city, facility.CityID)
	if result.Error != nil {
		panic(result.Error)
	}

	var facilityUsers []model.FacilityUser
	database.Db.Find(&facilityUsers, "FacilityId = ?", facility.ID)

	var ownerNames []string
	for _, fu := range facilityUsers {
		var user model.User
		result = database.Db.First(&user, fu.UserID)
		if result.Error == nil && user.RoleID == 2 {
			ownerNames = append(ownerNames, user.Name)
		}
	}

	exportDto := exportdtos.FacilityExportDto{
		ID:       facility.ID,
		Name:     facility.Name,
		Address:  facility.Address,
		City:     city.Name,
		Owners:   ownerNames,
		ImageURL: facility.ImageURL,
	}

	c.IndentedJSON(http.StatusOK, exportDto)
}

// func getAllFacilities(c *gin.Context) (model.Facility, error) {
// 	var facilities []model.Facility

// 	result := database.Db.First(&facilities)

// 	return facilities, nil
// }
