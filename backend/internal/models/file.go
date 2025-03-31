package models

type File struct {
	Name    string
	UserID  UserID
	GroupID StudyGroupID
}

type FileContext struct {
	Name string
	Data string
}
