package models

type ModuleID int64

type ModuleDetails struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Module struct {
	ID ModuleID `json:"id"`
	ModuleDetails
}
