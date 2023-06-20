package exportdtos

type FacilityExportDto struct {
	ID       int64    `gorm:"column:Id;primaryKey" json:"Id"`
	Name     string   `gorm:"column:Name;not null" json:"Name"`
	Address  string   `gorm:"column:Address;not null" json:"Address"`
	City     string   `json:"City"`
	Owners   []string `json:"Owners"`
	ImageURL string   `gorm:"column:ImageUrl;not null" json:"ImageUrl"`
}
