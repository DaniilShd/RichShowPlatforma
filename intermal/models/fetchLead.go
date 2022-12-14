package models

type FetchLead struct {
	RawLeads       int `json:"raw-leads"`
	ConfirmedLeads int `json:"confirmed-leads"`
	ArchiveLeads   int `json:"archive-leads"`
}
