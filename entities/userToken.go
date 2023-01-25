package entities

type UserToken struct {
	ID           uint   `json:"id"`
	UserId       *uint  `json:"user_id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
