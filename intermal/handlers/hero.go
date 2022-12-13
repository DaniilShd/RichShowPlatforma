package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/DaniilShd/RichShowPlatforma/intermal/forms"
	"github.com/DaniilShd/RichShowPlatforma/intermal/helpers"
	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
	"github.com/DaniilShd/RichShowPlatforma/intermal/render"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (m *Repository) HeroesAll(w http.ResponseWriter, r *http.Request) {
	heroes, err := m.DB.GetAllHeroes()
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-all-heroes.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "all-heroes.page.html"
	}

	data := make(map[string]interface{})
	data["heroes"] = heroes

	render.Template(w, r, template, &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) NewHero(w http.ResponseWriter, r *http.Request) {

	hero := models.Hero{}

	data := make(map[string]interface{})
	data["hero"] = hero

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-hero-new.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "hero-new.page.html"
	}

	render.Template(w, r, template, &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) NewPostHero(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm() // grab the multipart form
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	formdata := r.MultipartForm

	hero := models.Hero{}
	form := forms.New(r.PostForm)
	form.Required("name_hero")

	hero.Name = r.Form.Get("name_hero")
	hero.Description = r.Form.Get("description")
	hero.Gender, _ = strconv.Atoi(r.Form.Get("id_gender_hero"))
	hero.ClothingSizeMin, err = strconv.Atoi(r.Form.Get("clothing_size_min"))
	if err != nil {
		form.Errors.Add("clothing_size_min", "Необходимо ввести число!")
	}
	hero.ClothingSizeMax, err = strconv.Atoi(r.Form.Get("clothing_size_max"))
	if err != nil {
		form.Errors.Add("clothing_size_max", "Необходимо ввести число!")
	}

	err = r.ParseMultipartForm(8388608) // grab the multipart form
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	file := formdata.File["photo"]

	if len(file) == 0 {
		form.Errors.Add("photo", "Добавьте фото!")
	} else {
		photo := file[0].Filename
		hero.Photo = photo
	}

	if !form.Valid() {
		var template string
		switch m.App.Session.Get(r.Context(), "access_level") {
		case ADMIN_ACCESS_LEVEL:
			template = "admin-hero-new.page.html"
		case MANAGER_ACCESS_LEVEL:
			template = "hero-new.page.html"
		}

		data := make(map[string]interface{})
		data["hero"] = hero

		render.Template(w, r, template, &models.TemplateData{
			Data: data,
			Form: form,
		})
		return
	}

	uuidPhoto := uuid.New().String()
	suffix := strings.SplitAfter(file[0].Filename, ".")[1]

	out, err := os.Create("./static/img/heroes/" + uuidPhoto + "." + suffix)
	hero.Photo = uuidPhoto + "." + suffix
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	defer out.Close()

	fileOpen, err := file[0].Open()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	defer fileOpen.Close()

	_, err = io.Copy(out, fileOpen) // file not files[i] !
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.DB.InsertHero(&hero)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "/admin/heroes"
	case MANAGER_ACCESS_LEVEL:
		template = "/manager"
	}

	m.App.Session.Put(r.Context(), "flash", "Герой успешно добавлен")

	http.Redirect(w, r, template, http.StatusSeeOther)
}

func (m *Repository) DeleteHero(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "/admin/heroes"
	case MANAGER_ACCESS_LEVEL:
		template = "??????????????????"
	}

	err = m.DB.DeleteHeroByID(id)
	if err != nil {
		fmt.Println(err)
		m.App.Session.Put(r.Context(), "error", "Не удалось удалить героя")
		http.Redirect(w, r, template, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, template, http.StatusSeeOther)
}

func (m *Repository) ShowHero(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	hero, err := m.DB.GetHeroByID(id)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["hero"] = hero

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-hero-show.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "?????????????"
	}

	render.Template(w, r, template, &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) ChangeHero(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	hero, err := m.DB.GetHeroByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["hero"] = hero

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-hero-change.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "?????????????????"
	}

	render.Template(w, r, template, &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) ChangePostHero(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm() // grab the multipart form
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	formdata := r.MultipartForm

	hero := models.Hero{}
	form := forms.New(r.PostForm)
	form.Required("name_hero", "description")

	hero.ID, _ = strconv.Atoi(r.Form.Get("id"))
	hero.Name = r.Form.Get("name_hero")
	hero.Description = r.Form.Get("description")
	hero.ClothingSizeMin, err = strconv.Atoi(r.Form.Get("clothing_size_min"))
	if err != nil {
		form.Errors.Add("clothing_size_min", "Необходимо ввести число!")
	}
	hero.ClothingSizeMax, err = strconv.Atoi(r.Form.Get("clothing_size_max"))
	if err != nil {
		form.Errors.Add("clothing_size_max", "Необходимо ввести число!")
	}

	hero.Gender, _ = strconv.Atoi(r.Form.Get("id_gender_hero"))

	err = r.ParseMultipartForm(8388608) // grab the multipart form
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if !form.Valid() {
		var template string
		switch m.App.Session.Get(r.Context(), "access_level") {
		case ADMIN_ACCESS_LEVEL:
			template = "admin-hero-change.page.html"
		case MANAGER_ACCESS_LEVEL:
			template = "hero-new.page.html"
		}

		data := make(map[string]interface{})
		data["hero"] = hero

		render.Template(w, r, template, &models.TemplateData{
			Data: data,
			Form: form,
		})
		return
	}

	file := formdata.File["photo"]
	if len(file) != 0 {
		uuidPhoto := uuid.New().String()
		suffix := strings.SplitAfter(file[0].Filename, ".")[1]
		hero.Photo = uuidPhoto + "." + suffix
		out, err := os.Create("./static/img/heroes/" + uuidPhoto + "." + suffix)

		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		defer out.Close()

		fileOpen, err := file[0].Open()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		defer fileOpen.Close()

		_, err = io.Copy(out, fileOpen) // file not files[i] !
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
	} else {
		hero.Photo = ""
	}

	err = m.DB.UpdateHero(&hero)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "/admin/hero/" + strconv.Itoa(hero.ID)
	case MANAGER_ACCESS_LEVEL:
		template = "/manager" + strconv.Itoa(hero.ID)
	}

	m.App.Session.Put(r.Context(), "flash", "Герой успешно обновлен")

	http.Redirect(w, r, template, http.StatusSeeOther)
}
