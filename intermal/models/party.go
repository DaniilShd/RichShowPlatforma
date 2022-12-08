package models

type Party struct {
	ID          int
	Name        string
	Description string
	CheckList   CheckList
}
