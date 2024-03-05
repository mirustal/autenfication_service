package models


type AccessResponse struct {
	Access string `json:"access" bson:"_access,omitempty"`
	Refresh string `json:"refresh" bson:"_refresh, omitemty"`
}

