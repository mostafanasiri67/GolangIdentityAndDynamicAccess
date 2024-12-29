package entities

type User struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Firstname   string `json:"firstname,omitempty"`
	Lastname    string `json:"lastname,omitempty"`
	Password    string `json:"password,omitempty"`
	DateCreated string `json:"date_created,omitempty"`
	Role        string `json:"role,omitempty"`
	UserTokens  *[]UserToken
}
