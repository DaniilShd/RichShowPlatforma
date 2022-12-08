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

func (m *Repository) AllAnimation(w http.ResponseWriter, r *http.Request) {
	animationes, err := m.DB.GetAllAnimation()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["animation"] = animationes

	render.Template(w, r, "admin-all-animation.page.html", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ShowAnimation(w http.ResponseWriter, r *http.Request) {

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
	case "others":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_OTHER
	}

	res, err := m.DB.GetAnimationByID(id)
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
	data["animation"] = res
	data["check-lists"] = checkLists

	render.Template(w, r, "admin-animation-show.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) ShowPostAnimation(w http.ResponseWriter, r *http.Request) {

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

	animation, err := m.DB.GetAnimationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	animation.Name = r.Form.Get("name_animation")
	animation.Description = r.Form.Get("description")
	animation.CheckList.ID, _ = strconv.Atoi(r.Form.Get("id_check_list"))

	err = m.DB.UpdateAnimation(animation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")

	http.Redirect(w, r, "/admin/animation", http.StatusSeeOther)
}

// Show page of add new animation

func (m *Repository) NewAnimation(w http.ResponseWriter, r *http.Request) {

	checkLists, err := m.DB.GetAllCheckListsOfType(CHECK_LISTS_TYPE_OF_ANIMATION)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	animation := models.Animation{}

	data := make(map[string]interface{})
	data["check-lists"] = checkLists
	data["animation"] = animation

	render.Template(w, r, "admin-animation-new.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// Add new animation
func (m *Repository) NewPostAnimation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant't parse form!")
		http.Redirect(w, r, "admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	var animation models.Animation
	animation.Name = r.Form.Get("name_animation")
	animation.Description = r.Form.Get("description")
	animation.CheckList.ID, _ = strconv.Atoi(r.Form.Get("id_check_list"))

	form := forms.New(r.PostForm)
	form.Required("name_animation", "description", "id_check_list")

	if !form.Valid() {
		checkLists, err := m.DB.GetAllCheckList()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		data := make(map[string]interface{})
		data["check-lists"] = checkLists
		data["animation"] = animation
		render.Template(w, r, "admin-animation-new.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.InsertAnimation(&animation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "animation saved")

	http.Redirect(w, r, "/admin/animation", http.StatusSeeOther)
}

//Delete animation Deleteanimation

func (m *Repository) DeleteAnimation(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.DB.DeleteAnimationByID(id)
	if err != nil {
		http.Redirect(w, r, "/admin/animation", http.StatusSeeOther)
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/admin/animation", http.StatusSeeOther)
}
