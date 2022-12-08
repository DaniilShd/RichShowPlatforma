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

//Show-program ----------------------------------------------------------------------------------------------------

func (m *Repository) AllProgramShow(w http.ResponseWriter, r *http.Request) {
	programShows, err := m.DB.GetAllProgramShow()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["show-program"] = programShows

	render.Template(w, r, "admin-all-show-program.page.html", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ShowProgramShow(w http.ResponseWriter, r *http.Request) {

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

	res, err := m.DB.GetProgramShowByID(id)
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
	data["show-program"] = res
	data["check-lists"] = checkLists

	render.Template(w, r, "admin-show-program-show.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) ShowPostProgramShow(w http.ResponseWriter, r *http.Request) {

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

	programShow, err := m.DB.GetProgramShowByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	programShow.Name = r.Form.Get("name_show")
	programShow.Description = r.Form.Get("description")
	programShow.CheckList.ID, _ = strconv.Atoi(r.Form.Get("id_check_list"))

	err = m.DB.UpdateProgramShow(programShow)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")

	http.Redirect(w, r, "/admin/show-program", http.StatusSeeOther)
}

// Show page of add new show-program

func (m *Repository) NewProgramShow(w http.ResponseWriter, r *http.Request) {

	checkLists, err := m.DB.GetAllCheckListsOfType(CHECK_LISTS_TYPE_OF_SHOW)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	programShow := models.ProgramShow{}

	data := make(map[string]interface{})
	data["check-lists"] = checkLists
	data["show-program"] = programShow

	render.Template(w, r, "admin-show-program-new.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// Add new show-program
func (m *Repository) NewPostProgramShow(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant't parse form!")
		http.Redirect(w, r, "admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	var programShow models.ProgramShow
	programShow.Name = r.Form.Get("name_show")
	programShow.Description = r.Form.Get("description")
	programShow.CheckList.ID, err = strconv.Atoi(r.Form.Get("id_check_list"))
	if err != nil {

		fmt.Println(err)
	}

	form := forms.New(r.PostForm)
	form.Required("name_show", "description", "id_check_list")

	if !form.Valid() {
		checkLists, err := m.DB.GetAllCheckList()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		data := make(map[string]interface{})
		data["check-lists"] = checkLists
		data["show-program"] = programShow
		render.Template(w, r, "admin-show-program-new.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.InsertProgramShow(&programShow)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "show-program saved")

	http.Redirect(w, r, "/admin/show-program", http.StatusSeeOther)
}

//Delete show-program DeleteprogramShow

func (m *Repository) DeleteProgramShow(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.DB.DeleteProgramShowByID(id)
	if err != nil {
		http.Redirect(w, r, "/admin/show-program", http.StatusSeeOther)
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/admin/show-program", http.StatusSeeOther)
}
