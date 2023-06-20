package importdtos

type FacilityImportDto struct {
	Name      string `json:"Name"`
	Address   string `json:"Address"`
	City      string `json:"City"`
	ImageURL  string `json:"ImageUrl"`
	CreatorId string `json:"CreatorId"`
}
