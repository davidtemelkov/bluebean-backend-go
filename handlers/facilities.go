package handlers

import (
	"example/bluebean-go/database"
	exportdtos "example/bluebean-go/dtos/export"
	importdtos "example/bluebean-go/dtos/import"
	"example/bluebean-go/model"
	"example/bluebean-go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func mapFacilityToExportDto(f model.Facility, facilityUsers []model.FacilityUser) exportdtos.FacilityExportDto {
	var city model.City
	result := database.Db.First(&city, f.CityID)
	if result.Error != nil {
		panic(result.Error)
	}

	var ownerNames []string
	for _, fu := range facilityUsers {
		var user model.User
		result = database.Db.First(&user, fu.UserID)
		if result.Error == nil && user.RoleID == 2 {
			ownerNames = append(ownerNames, user.Name)
		}
	}

	exportDto := exportdtos.FacilityExportDto{
		ID:       f.ID,
		Name:     f.Name,
		Address:  f.Address,
		City:     city.Name,
		Owners:   ownerNames,
		ImageURL: f.ImageURL,
	}

	return exportDto
}

func GetAllFacilities(c *gin.Context) {
	var facilities []model.Facility

	result := database.Db.Find(&facilities)
	if result.Error != nil {
		panic(result.Error)
	}

	var exportDtos []exportdtos.FacilityExportDto

	for _, f := range facilities {
		var city model.City
		result = database.Db.First(&city, f.CityID)
		if result.Error != nil {
			panic(result.Error)
		}

		var facilityUsers []model.FacilityUser
		database.Db.Find(&facilityUsers, "FacilityId = ?", f.ID)

		exportDto := mapFacilityToExportDto(f, facilityUsers)
		exportDtos = append(exportDtos, exportDto)
	}

	c.IndentedJSON(http.StatusOK, exportDtos)
}

func GetOneFacility(c *gin.Context) {
	id := c.Param("facilityId")

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

	exportDto := mapFacilityToExportDto(facility, facilityUsers)

	c.IndentedJSON(http.StatusOK, exportDto)
}

func GetAllFacilitiesForUser(c *gin.Context) {
	userId := c.Param("userId")

	var user model.User
	result := database.Db.First(&user, userId)
	if result.Error != nil {
		panic(result.Error)
	}

	var facilityUsers []model.FacilityUser
	database.Db.Find(&facilityUsers, "UserId = ?", user.ID)

	var facilityIDs []int64
	for _, fu := range facilityUsers {
		facilityIDs = append(facilityIDs, fu.FacilityID)
	}

	var facilities []model.Facility
	database.Db.Find(&facilities, facilityIDs)

	var exportDtos []exportdtos.FacilityExportDto

	for _, f := range facilities {
		var facilityUsers []model.FacilityUser
		database.Db.Find(&facilityUsers, "FacilityId = ?", f.ID)

		exportDto := mapFacilityToExportDto(f, facilityUsers)
		exportDtos = append(exportDtos, exportDto)
	}

	c.IndentedJSON(http.StatusOK, exportDtos)
}

func CreateFacility(c *gin.Context) {
	var importDto importdtos.FacilityImportDto

	if err := c.BindJSON(&importDto); err != nil {
		panic("Wrong json format.")
	}

	var city model.City
	result := database.Db.Where("name = ?", importDto.City).First(&city)
	if result.Error != nil {
		panic(result.Error)
	}

	creatorID, err := strconv.ParseInt(importDto.CreatorId, 10, 64)
	if err != nil {
		panic(err)
	}

	imageUrl, err := utils.UploadFile(importDto.ImageURL, importDto.Name, "Photo")

	newFacility := model.Facility{
		Name:      importDto.Name,
		Address:   importDto.Address,
		CityID:    city.ID,
		CreatorID: creatorID,
		ImageURL:  imageUrl,
	}

	if err := database.Db.Create(&newFacility).Error; err != nil {
		panic("Error creating facility.")
	}

	c.IndentedJSON(http.StatusOK, newFacility)
}
