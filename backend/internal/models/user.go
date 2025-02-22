package models

type UserID string
type User struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Modules []Module `json:"modules"`
}
