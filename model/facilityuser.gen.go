// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameFacilityUser = "FacilityUser"

// FacilityUser mapped from table <FacilityUser>
type FacilityUser struct {
	FacilityID int64 `gorm:"column:FacilityId;primaryKey" json:"FacilityId"`
	UserID     int64 `gorm:"column:UserId;primaryKey" json:"UserId"`
}

// TableName FacilityUser's table name
func (*FacilityUser) TableName() string {
	return TableNameFacilityUser
}