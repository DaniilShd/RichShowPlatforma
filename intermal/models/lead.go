package models

import "time"

type Program struct {
	ID       int
	Name     string
	Duration string
}

type Artist struct {
	ID          int
	Name        string
	PhoneNumber int
	Photo       string
}

type Hero struct {
	ID     int
	Name   string
	Gender string
	Artist Artist
}

type Child struct {
	ID             int
	Name           string
	Gender         int
	Artist         Artist
	DateOfBirthDay time.Time
}

type Lead struct {
	ID                   int       `db:"id_lead"`
	FirstNameClient      string    `db:"first_name"`
	LastNameClient       string    `db:"last_name"`
	AmountOfChildren     int       `db:"amount_of_children"`
	AverageAgeOfChildren int       `db:"average_age_of_children"`
	Address              string    `db:"address"`
	Date                 time.Time `db:"date"`
	ActiveLead           bool      `db:"active_lead"`
	CheckArtists         bool      `db:"check_artists"`
	Confirmed            bool      `db:"confirmed"`
	CheckAssistants      bool      `db:"check_assistants"`
	PhoneNumber          int
	Telegram             string
	Programs             []Program
	Heroes               []Hero
	Child                Child
}
