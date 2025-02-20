package models

type ModuleID int64 // I'll leave this for now
type Module struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
