package models

type CheckList struct {
	ID           int
	Name         string
	Description  string
	TypeOfList   int
	NameOfPoints []string
}
