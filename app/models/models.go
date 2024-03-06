package models

type AccessResponse struct {

	Access string `json:"access" bson:"_access,omitempty"`
	Refresh string `json:"refresh" bson:"_refresh, omitemty"`
}

type ResponseDB struct {
	GUID string	`json: "guid" bson:"_guid"`
	RefreshToken []byte `json:"refreshtoken" bson:"_refreshtoken,omitempty`
}

type User struct {
	GUID string	`json: "guid" bson:"_guid"`
	TokenUser AccessResponse `json: "token" bson:"_token"`
}

type UserCookie struct {
	AccessToken string `cookie:"accesstoken"`
	RefreshToken string `cookie:"refreshtoken"`
}