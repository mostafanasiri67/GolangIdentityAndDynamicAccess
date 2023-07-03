package models

type Error struct {
	ResponseCode      int    `json:"rc"`
	Message           string `json:"message"`
	Detail            string `json:"detail"`
	ExternalReference string `json:"ext_ref"`
}

type ApiConfig struct {
	LoginPath             string `json:"loginPath"`
	RegisterPath          string `json:"registerPath"`
	LogoutPath            string `json:"logoutPath"`
	RefreshTokenPath      string `json:"refreshTokenPath"`
	AccessTokenObjectKey  string `json:"accessTokenObjectKey"`
	RefreshTokenObjectKey string `json:"refreshTokenObjectKey"`
	AdminRoleName         string `json:"adminRoleName"`
}
