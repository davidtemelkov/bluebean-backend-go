// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameRole = "Role"

// Role mapped from table <Role>
type Role struct {
	ID   int64  `gorm:"column:Id;primaryKey" json:"Id"`
	Name string `gorm:"column:Name;not null" json:"Name"`
}

// TableName Role's table name
func (*Role) TableName() string {
	return TableNameRole
}