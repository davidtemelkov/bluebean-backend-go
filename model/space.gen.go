// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameSpace = "Space"

// Space mapped from table <Space>
type Space struct {
	ID         int64  `gorm:"column:Id;primaryKey" json:"Id"`
	Name       string `gorm:"column:Name;not null" json:"Name"`
	Location   string `gorm:"column:Location;not null" json:"Location"`
	SchemaURL  string `gorm:"column:SchemaUrl;not null" json:"SchemaUrl"`
	FacilityID int64  `gorm:"column:FacilityId;not null" json:"FacilityId"`
}

// TableName Space's table name
func (*Space) TableName() string {
	return TableNameSpace
}
