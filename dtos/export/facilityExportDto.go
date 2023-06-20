package exportdtos

type FacilityExportDto struct {
	ID       int64    `json:"Id"`
	Name     string   `json:"Name"`
	Address  string   `json:"Address"`
	City     string   `json:"City"`
	Owners   []string `json:"Owners"`
	ImageURL string   `json:"ImageUrl"`
}
