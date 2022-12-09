package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/DaniilShd/RichShowPlatforma/intermal/forms"
	"github.com/DaniilShd/RichShowPlatforma/intermal/helpers"
	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
	"github.com/DaniilShd/RichShowPlatforma/intermal/render"
	"github.com/go-chi/chi"
)

// func (m *Repository) Doashboard(w http.ResponseWriter, r *http.Request) {
// 	render.Template(w, r, "admin-dashboard.page.html", &models.TemplateData{})
// }

//party ----------------------------------------------------------------------------------------------------

func (m *Repository) AllParty(w http.ResponseWriter, r *http.Request) {
	partys, err := m.DB.GetAllParty()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["party"] = partys

	render.Template(w, r, "admin-all-party.page.html", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ShowParty(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[3])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[2]

	var typeOfCheckList int
	switch src {
	case "show-program":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_SHOW
	case "master-class":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_MASTER_CLASS
	case "animations":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_ANIMATION
	case "party":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_PARTIES_AND_QUESTS
	case "other":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_OTHER
	}

	res, err := m.DB.GetPartyByID(id)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	checkLists, err := m.DB.GetAllCheckListsOfType(typeOfCheckList)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["party"] = res
	data["check-lists"] = checkLists

	render.Template(w, r, "admin-party-show.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) ShowPostParty(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant't parse form!")
		http.Redirect(w, r, "admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[3])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	Party, err := m.DB.GetPartyByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	Party.Name = r.Form.Get("name_party_quest")
	Party.Description = r.Form.Get("description")
	Party.CheckList.ID, _ = strconv.Atoi(r.Form.Get("id_check_list"))

	err = m.DB.UpdateParty(Party)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")

	http.Redirect(w, r, "/admin/party", http.StatusSeeOther)
}

// Show page of add new party

func (m *Repository) NewParty(w http.ResponseWriter, r *http.Request) {

	checkLists, err := m.DB.GetAllCheckListsOfType(CHECK_LISTS_TYPE_OF_PARTIES_AND_QUESTS)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	party := models.Party{}

	data := make(map[string]interface{})
	data["check-lists"] = checkLists
	data["party"] = party

	render.Template(w, r, "admin-party-new.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// Add new party
func (m *Repository) NewPostParty(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant't parse form!")
		http.Redirect(w, r, "admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	var Party models.Party
	Party.Name = r.Form.Get("name_party_quest")
	Party.Description = r.Form.Get("description")
	Party.CheckList.ID, err = strconv.Atoi(r.Form.Get("id_check_list"))
	if err != nil {

		fmt.Println(err)
	}

	form := forms.New(r.PostForm)
	form.Required("name_party_quest", "description", "id_check_list")

	if !form.Valid() {
		checkLists, err := m.DB.GetAllCheckList()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		data := make(map[string]interface{})
		data["check-lists"] = checkLists
		data["party"] = Party
		render.Template(w, r, "admin-party-new.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.InsertParty(&Party)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "party saved")

	http.Redirect(w, r, "/admin/party", http.StatusSeeOther)
}

//Delete party DeleteParty

func (m *Repository) DeleteParty(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.DB.DeletePartyByID(id)
	if err != nil {
		http.Redirect(w, r, "/admin/party", http.StatusSeeOther)
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/admin/party", http.StatusSeeOther)
}