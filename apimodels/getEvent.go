package apimodels

type GetEvent struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	StartDate          string `json:"startDate"`
	EndDate            string `json:"endDate"`
	VolunteersRequired int    `json:"volunteersRequired"`
	Location           string `json:"location"`
}
