package models

type StoreItem struct {
	ID            int
	Name          string
	CurrentAmount int
	Dimension     string
	MinAmount     int
	Description   string
}
