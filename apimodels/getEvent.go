package apimodels

import "time"

type GetEvent struct {
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	StartDate          time.Time `json:"startDate"`
	EndDate            time.Time `json:"endDate"`
	VolunteersRequired int       `json:"volunteersRequired"`
	Location           string    `json:"location"`
}
