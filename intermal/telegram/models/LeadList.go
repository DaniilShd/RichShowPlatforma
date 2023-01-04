package models

import (
	"time"

	modelsFromApp "github.com/DaniilShd/RichShowPlatforma/intermal/models"
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
	Command             string                   `json:"command"`
	LeadID              int                      `json:"id_lead"`
	ChatID              int64                    `json:"id_chat"`
	ResponseLeadFromApp chan *modelsFromApp.Lead `json:"-"`
}
