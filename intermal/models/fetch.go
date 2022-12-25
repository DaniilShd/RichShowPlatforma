package models

type FetchLead struct {
	RawLeads       int `json:"raw-leads"`
	ConfirmedLeads int `json:"confirmed-leads"`
	ArchiveLeads   int `json:"archive-leads"`
}

type FetchStoreOrder struct {
	NewOrder       int `json:"new-order"`
	CompletedOrder int `json:"completed-order"`
	DestroyOrder   int `json:"destroy-order"`
}
