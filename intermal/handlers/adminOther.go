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

//Other ----------------------------------------------------------------------------------------------------

func (m *Repository) AllOther(w http.ResponseWriter, r *http.Request) {
	others, err := m.DB.GetAllOther()
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["other"] = others

	render.Template(w, r, "admin-all-other.page.html", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ShowOther(w http.ResponseWriter, r *http.Request) {

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

	res, err := m.DB.GetOtherByID(id)
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
	data["other"] = res
	data["check-lists"] = checkLists

	render.Template(w, r, "admin-other-show.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) ShowPostOther(w http.ResponseWriter, r *http.Request) {

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

	Other, err := m.DB.GetOtherByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	Other.Name = r.Form.Get("name_other")
	Other.Description = r.Form.Get("description")
	Other.CheckList.ID, _ = strconv.Atoi(r.Form.Get("id_check_list"))

	err = m.DB.UpdateOther(Other)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")

	http.Redirect(w, r, "/admin/other", http.StatusSeeOther)
}

// Show page of add new Other

func (m *Repository) NewOther(w http.ResponseWriter, r *http.Request) {

	checkLists, err := m.DB.GetAllCheckListsOfType(CHECK_LISTS_TYPE_OF_PARTIES_AND_QUESTS)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	other := models.Other{}

	data := make(map[string]interface{})
	data["check-lists"] = checkLists
	data["other"] = other

	render.Template(w, r, "admin-other-new.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// Add new Other
func (m *Repository) NewPostOther(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant't parse form!")
		http.Redirect(w, r, "admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	var Other models.Other
	Other.Name = r.Form.Get("name_other")
	Other.Description = r.Form.Get("description")
	Other.CheckList.ID, err = strconv.Atoi(r.Form.Get("id_check_list"))
	if err != nil {

		fmt.Println(err)
	}

	form := forms.New(r.PostForm)
	form.Required("name_other", "description", "id_check_list")

	if !form.Valid() {
		checkLists, err := m.DB.GetAllCheckList()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		data := make(map[string]interface{})
		data["check-lists"] = checkLists
		data["other"] = Other
		render.Template(w, r, "admin-other-new.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.InsertOther(&Other)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "other saved")

	http.Redirect(w, r, "/admin/other", http.StatusSeeOther)
}

//Delete Other DeleteOther

func (m *Repository) DeleteOther(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.DB.DeleteOtherByID(id)
	if err != nil {
		http.Redirect(w, r, "/admin/other", http.StatusSeeOther)
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/admin/other", http.StatusSeeOther)
}
