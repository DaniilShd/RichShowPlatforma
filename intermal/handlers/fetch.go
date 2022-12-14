package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DaniilShd/RichShowPlatforma/intermal/helpers"
	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

func (m *Repository) FetchLead(w http.ResponseWriter, r *http.Request) {

	var fetchResult models.FetchLead
	var err error

	fetchResult.RawLeads, err = m.DB.GetCountRawLeads()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	fetchResult.ArchiveLeads, err = m.DB.GetCountArchiveLeads()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	fetchResult.ConfirmedLeads, err = m.DB.GetCountConfirmedLeads()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	js, err := json.Marshal(fetchResult)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(js)
}
