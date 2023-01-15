package models

import "time"

type StoreLead struct {
	ID               int
	Name             string
	ProgramType      string
	LeadID           int
	CheckListID      int
	Date             time.Time
	Time             time.Time
	Photo            string
	StoreDescription string
	LeadDescription  string
	AmountOfChilds   int
	Canceled         bool
	Completed        bool
}
