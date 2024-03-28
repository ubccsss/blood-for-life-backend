package apimodels

type UnregisterEvent struct {
	UserID  int `json:"userID"`
	EventID int `json:"eventID"`
}

/*

Logic for deregistering user

1. Get context of the user and their ID
2. Get context of the event and its ID 

4. Check if user is on signup list and remove 
	5. If not, check if user is on standby list and remove 

*/

