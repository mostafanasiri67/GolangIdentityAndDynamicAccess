package models

import "goLang/entities"

type Login struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}
type RefreshToken struct {
	RefreshToken string `json:"RefreshToken"`
}
type Permission struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path,omitempty"`
	UserId uint   `json:"userId,omitempty"`
}
type UserPermission struct {
	Permission     string `json:"Permission,omitempty"`
	UserPermission []entities.Permission
}
type User struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	Firstname    string `json:"firstname,omitempty"`
	Lastname     string `json:"lastname,omitempty"`
	Password     string `json:"password,omitempty"`
	DateCreated  string `json:"date_created,omitempty"`
	Role         string `json:"role,omitempty"`
	AccessToken  string `json:"AccessToken,omitempty"`
	RefreshToken string `json:"RefreshToken,omitempty"`
	Permission   string `json:"Permission,omitempty"`
}

type UserToken struct {
	ID           uint   `json:"id"`
	UserId       uint   `json:"user_id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type ChangePassword struct {
	OldPassword     string `json:"oldPassword"`
	NewPassword     string `json:"newPassword,omitempty"`
	ConfirmPassword string `json:"confirmPassword,omitempty"`
}
