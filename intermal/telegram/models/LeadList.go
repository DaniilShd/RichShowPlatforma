package models

import (
	"time"
)

type LeadList struct {
	ID         int       `db:"id_lead"`
	Address    string    `db:"address"`
	Date       time.Time `db:"date"`
	Time       time.Time `db:"time"`
	NameHeroes []string
	ArtistID   int
}

type RequestFromChat struct {
	Command             string           `json:"command"`
	LeadID              int              `json:"id_lead"`
	ChatID              int64            `json:"id_chat"`
	CheckListID         int              `json:"-"`
	ResponseLeadFromApp chan interface{} `json:"-"`
	StoreOrderID        []int            `json:"-"`
}
