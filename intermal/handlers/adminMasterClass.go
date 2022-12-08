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

func (m *Repository) Doashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.html", &models.TemplateData{})
}

//Master-class ----------------------------------------------------------------------------------------------------

func (m *Repository) AllMasterClass(w http.ResponseWriter, r *http.Request) {
	masterClasses, err := m.DB.GetAllMasterClass()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["master-class"] = masterClasses

	render.Template(w, r, "admin-all-master-class.page.html", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ShowMasterClass(w http.ResponseWriter, r *http.Request) {

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
	case "animation":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_ANIMATION
	case "parties":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_PARTIES_AND_QUESTS
	case "other":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_OTHER
	}

	res, err := m.DB.GetMasterClassByID(id)
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
	data["master-class"] = res
	data["check-lists"] = checkLists

	render.Template(w, r, "admin-master-class-show.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) ShowPostMasterClass(w http.ResponseWriter, r *http.Request) {

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

	masterClass, err := m.DB.GetMasterClassByID(id)

	masterClass.Name = r.Form.Get("name_master_class")
	masterClass.Description = r.Form.Get("description")
	masterClass.CheckList.ID, _ = strconv.Atoi(r.Form.Get("id_check_list"))

	err = m.DB.UpdateMasterClass(masterClass)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")

	http.Redirect(w, r, "/admin/master-class", http.StatusSeeOther)
}

// Show page of add new master-class

func (m *Repository) NewMasterClass(w http.ResponseWriter, r *http.Request) {

	checkLists, err := m.DB.GetAllCheckListsOfType(CHECK_LISTS_TYPE_OF_MASTER_CLASS)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	masterClass := models.MasterClass{}

	data := make(map[string]interface{})
	data["check-lists"] = checkLists
	data["master-class"] = masterClass

	render.Template(w, r, "admin-master-class-new.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// Add new master-class
func (m *Repository) NewPostMasterClass(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant't parse form!")
		http.Redirect(w, r, "admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	var masterClass models.MasterClass
	masterClass.Name = r.Form.Get("name_master_class")
	masterClass.Description = r.Form.Get("description")
	masterClass.CheckList.ID, _ = strconv.Atoi(r.Form.Get("id_check_list"))

	form := forms.New(r.PostForm)
	form.Required("name_master_class", "description", "id_check_list")

	if !form.Valid() {
		checkLists, err := m.DB.GetAllCheckList()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		data := make(map[string]interface{})
		data["check-lists"] = checkLists
		data["master-class"] = masterClass
		render.Template(w, r, "admin-master-class-new.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.InsertMasterClass(&masterClass)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Master-class saved")

	http.Redirect(w, r, "/admin/master-class", http.StatusSeeOther)
}

//Delete master-class DeleteMasterClass

func (m *Repository) DeleteMasterClass(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.DB.DeleteMasterClassByID(id)
	if err != nil {
		http.Redirect(w, r, "/admin/master-class", http.StatusSeeOther)
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/admin/master-class", http.StatusSeeOther)
}

//Leads ----------------------------------------------------------------------------------------------------------------

func (m *Repository) LeadsCalendar(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-leads-calendar.page.html", &models.TemplateData{})
}
