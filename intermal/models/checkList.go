package models

type Item struct {
	ID             int
	Name           string
	AmountItemOnce float64
	Dimension      string
}

type CheckList struct {
	ID           int
	Name         string
	Description  string
	TypeOfList   int
	NameOfPoints []string
	Items        []Item
	Duration     int
}
