package entities

type Permission struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path,omitempty"`
	UserId uint `json:"userId,omitempty"`
}
