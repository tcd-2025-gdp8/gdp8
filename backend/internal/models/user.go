package models

type UserID string

type UserDetails struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type User struct {
	ID UserID `json:"id"`
	UserDetails
	Modules []Module `json:"modules"`
}
