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

func (m *Repository) AnimatorsAll(w http.ResponseWriter, r *http.Request) {
	animators, err := m.DB.GetAllArtists()
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-all-animators.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "all-animators.page.html"
	}

	data := make(map[string]interface{})
	data["animators"] = animators

	render.Template(w, r, template, &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) NewAnimator(w http.ResponseWriter, r *http.Request) {

	animator := models.Artist{}

	data := make(map[string]interface{})
	data["animator"] = animator

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-animator-new.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "animator-new.page.html"
	}

	render.Template(w, r, template, &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) NewPostAnimator(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm() // grab the multipart form
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	formdata := r.MultipartForm

	animator := models.Artist{}
	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "phone_number")
	form.MinLength("phone_number", 17, r)

	animator.FirstName = r.Form.Get("first_name")
	animator.LastName = r.Form.Get("last_name")
	animator.Description = r.Form.Get("description")
	animator.Telegram = r.Form.Get("telegram_artist")
	animator.VK = r.Form.Get("vk")

	//преобразуем номер телефона к виду 9200261804
	phoneNumber := r.Form.Get("phone_number")
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "-", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "(", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, ")", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "+", "")
	if len(phoneNumber) > 1 {
		phoneNumber = phoneNumber[1:]
	}
	animator.PhoneNumber = phoneNumber
	animator.Growth, err = strconv.Atoi(r.Form.Get("growth"))
	if err != nil {
		form.Errors.Add("growth", "Необходимо ввести число!")
	}
	animator.ShoeSize, err = strconv.Atoi(r.Form.Get("shoe_size"))
	if err != nil {
		form.Errors.Add("shoe_size", "Необходимо ввести число!")
	}
	animator.ClothingSize, err = strconv.Atoi(r.Form.Get("clothing_size"))
	if err != nil {
		form.Errors.Add("clothing_size", "Необходимо ввести число!")
	}
	animator.Gender, _ = strconv.Atoi(r.Form.Get("id_gender_hero"))

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
		animator.Photo = photo
	}

	if !form.Valid() {
		animator.PhoneNumber = helpers.ConvertNumberPhone(animator.PhoneNumber)
		var template string
		switch m.App.Session.Get(r.Context(), "access_level") {
		case ADMIN_ACCESS_LEVEL:
			template = "admin-animator-new.page.html"
		case MANAGER_ACCESS_LEVEL:
			template = "animator-new.page.html"
		}

		data := make(map[string]interface{})
		data["animator"] = animator

		render.Template(w, r, template, &models.TemplateData{
			Data: data,
			Form: form,
		})
		return
	}

	uuidPhoto := uuid.New().String()
	suffix := strings.SplitAfter(file[0].Filename, ".")[1]

	out, err := os.Create("./static/img/animators/" + uuidPhoto + "." + suffix)
	animator.Photo = uuidPhoto + "." + suffix
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

	err = m.DB.InsertArtist(&animator)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "/admin/animators"
	case MANAGER_ACCESS_LEVEL:
		template = "/manager"
	}

	m.App.Session.Put(r.Context(), "flash", "Аниматор успешно добавлен")

	http.Redirect(w, r, template, http.StatusSeeOther)
}

func (m *Repository) DeleteAnimator(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "/admin/animators"
	case MANAGER_ACCESS_LEVEL:
		template = "/manager"
	}

	err = m.DB.DeleteArtistByID(id)
	if err != nil {
		fmt.Println(err)
		m.App.Session.Put(r.Context(), "error", "Не удалось удалить аниматора")
		http.Redirect(w, r, template, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, template, http.StatusSeeOther)
}

func (m *Repository) ShowAnimator(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	animator, err := m.DB.GetArtistByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	animator.PhoneNumber = helpers.ConvertNumberPhone(animator.PhoneNumber)

	data := make(map[string]interface{})
	data["animator"] = animator

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-animator-show.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "?????????????"
	}

	render.Template(w, r, template, &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) ChangeAnimator(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	animator, err := m.DB.GetArtistByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	animator.PhoneNumber = helpers.ConvertNumberPhone(animator.PhoneNumber)

	data := make(map[string]interface{})
	data["animator"] = animator

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-animator-change.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "?????????????????"
	}

	render.Template(w, r, template, &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) ChangePostAnimator(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm() // grab the multipart form
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	formdata := r.MultipartForm

	animator := models.Artist{}
	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "phone_number")
	form.MinLength("phone_number", 17, r)

	animator.ID, _ = strconv.Atoi(r.Form.Get("id"))
	animator.FirstName = r.Form.Get("first_name")
	animator.LastName = r.Form.Get("last_name")
	animator.Description = r.Form.Get("description")
	animator.Telegram = r.Form.Get("telegram_artist")
	animator.VK = r.Form.Get("vk")

	//преобразуем номер телефона к виду 9200261804
	phoneNumber := r.Form.Get("phone_number")
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "-", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "(", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, ")", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "+", "")
	if len(phoneNumber) > 1 {
		phoneNumber = phoneNumber[1:]
	}
	animator.PhoneNumber = phoneNumber
	animator.Growth, err = strconv.Atoi(r.Form.Get("growth"))
	if err != nil {
		form.Errors.Add("growth", "Необходимо ввести число!")
	}
	animator.ShoeSize, err = strconv.Atoi(r.Form.Get("shoe_size"))
	if err != nil {
		form.Errors.Add("shoe_size", "Необходимо ввести число!")
	}
	animator.ClothingSize, err = strconv.Atoi(r.Form.Get("clothing_size"))
	if err != nil {
		form.Errors.Add("clothing_size", "Необходимо ввести число!")
	}
	animator.Gender, _ = strconv.Atoi(r.Form.Get("id_gender_hero"))

	err = r.ParseMultipartForm(8388608) // grab the multipart form
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if !form.Valid() {
		animator.PhoneNumber = helpers.ConvertNumberPhone(animator.PhoneNumber)
		var template string
		switch m.App.Session.Get(r.Context(), "access_level") {
		case ADMIN_ACCESS_LEVEL:
			template = "admin-animator-change.page.html"
		case MANAGER_ACCESS_LEVEL:
			template = "animator-new.page.html"
		}

		data := make(map[string]interface{})
		data["animator"] = animator

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
		animator.Photo = uuidPhoto + "." + suffix
		out, err := os.Create("./static/img/animators/" + uuidPhoto + "." + suffix)

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
		animator.Photo = ""
	}

	err = m.DB.UpdateArtist(&animator)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "/admin/animator/" + strconv.Itoa(animator.ID)
	case MANAGER_ACCESS_LEVEL:
		template = "/manager" + strconv.Itoa(animator.ID)
	}

	m.App.Session.Put(r.Context(), "flash", "Аниматор успешно обновлен")

	http.Redirect(w, r, template, http.StatusSeeOther)
}
