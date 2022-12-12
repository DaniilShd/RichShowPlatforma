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

func (m *Repository) AssistantsAll(w http.ResponseWriter, r *http.Request) {
	assistants, err := m.DB.GetAllAssistants()
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-all-assistants.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "all-assistants.page.html"
	}

	data := make(map[string]interface{})
	data["assistants"] = assistants

	render.Template(w, r, template, &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) NewAssistant(w http.ResponseWriter, r *http.Request) {

	assistant := models.Assistant{}

	data := make(map[string]interface{})
	data["assistant"] = assistant

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-assistant-new.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "assistant-new.page.html"
	}

	render.Template(w, r, template, &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) NewPostAssistant(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm() // grab the multipart form
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	formdata := r.MultipartForm

	assistant := models.Assistant{}
	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "phone_number")
	form.MinLength("phone_number", 17, r)

	assistant.FirstName = r.Form.Get("first_name")
	assistant.LastName = r.Form.Get("last_name")
	assistant.Description = r.Form.Get("description")
	assistant.Telegram = r.Form.Get("telegram_assistant")
	assistant.VK = r.Form.Get("vk")

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
	assistant.PhoneNumber = phoneNumber
	assistant.Gender, _ = strconv.Atoi(r.Form.Get("id_gender_type"))

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
		assistant.Photo = photo
	}

	if !form.Valid() {
		assistant.PhoneNumber = helpers.ConvertNumberPhone(assistant.PhoneNumber)
		var template string
		switch m.App.Session.Get(r.Context(), "access_level") {
		case ADMIN_ACCESS_LEVEL:
			template = "admin-assistant-new.page.html"
		case MANAGER_ACCESS_LEVEL:
			template = "assistant-new.page.html"
		}

		data := make(map[string]interface{})
		data["assistant"] = assistant

		render.Template(w, r, template, &models.TemplateData{
			Data: data,
			Form: form,
		})
		return
	}

	uuidPhoto := uuid.New().String()
	suffix := strings.SplitAfter(file[0].Filename, ".")[1]

	out, err := os.Create("./static/img/assistants/" + uuidPhoto + "." + suffix)
	assistant.Photo = uuidPhoto + "." + suffix
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

	err = m.DB.InsertAssistant(&assistant)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "/admin/assistants"
	case MANAGER_ACCESS_LEVEL:
		template = "/manager"
	}

	m.App.Session.Put(r.Context(), "flash", "Аниматор успешно добавлен")

	http.Redirect(w, r, template, http.StatusSeeOther)
}

func (m *Repository) DeleteAssistant(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "/admin/assistants"
	case MANAGER_ACCESS_LEVEL:
		template = "/manager"
	}

	err = m.DB.DeleteAssistantByID(id)
	if err != nil {
		fmt.Println(err)
		m.App.Session.Put(r.Context(), "error", "Не удалось удалить ассистента")
		http.Redirect(w, r, template, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, template, http.StatusSeeOther)
}

func (m *Repository) ShowAssistant(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	assistant, err := m.DB.GetAssistantByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	assistant.PhoneNumber = helpers.ConvertNumberPhone(assistant.PhoneNumber)

	data := make(map[string]interface{})
	data["assistant"] = assistant

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-assistant-show.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "?????????????"
	}

	render.Template(w, r, template, &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) ChangeAssistant(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	assistant, err := m.DB.GetAssistantByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	assistant.PhoneNumber = helpers.ConvertNumberPhone(assistant.PhoneNumber)

	data := make(map[string]interface{})
	data["assistant"] = assistant

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-assistant-change.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "?????????????????"
	}

	render.Template(w, r, template, &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) ChangePostAssistant(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm() // grab the multipart form
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	formdata := r.MultipartForm

	assistant := models.Assistant{}
	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "phone_number")
	form.MinLength("phone_number", 17, r)

	assistant.ID, _ = strconv.Atoi(r.Form.Get("id"))
	assistant.FirstName = r.Form.Get("first_name")
	assistant.LastName = r.Form.Get("last_name")
	assistant.Description = r.Form.Get("description")
	assistant.Telegram = r.Form.Get("telegram_assistant")
	assistant.VK = r.Form.Get("vk")

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
	assistant.PhoneNumber = phoneNumber
	assistant.Gender, _ = strconv.Atoi(r.Form.Get("id_gender_type"))

	err = r.ParseMultipartForm(8388608) // grab the multipart form
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if !form.Valid() {
		assistant.PhoneNumber = helpers.ConvertNumberPhone(assistant.PhoneNumber)
		var template string
		switch m.App.Session.Get(r.Context(), "access_level") {
		case ADMIN_ACCESS_LEVEL:
			template = "admin-assistant-change.page.html"
		case MANAGER_ACCESS_LEVEL:
			template = "assistant-new.page.html"
		}

		data := make(map[string]interface{})
		data["assistant"] = assistant

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
		assistant.Photo = uuidPhoto + "." + suffix
		out, err := os.Create("./static/img/assistants/" + uuidPhoto + "." + suffix)

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
		assistant.Photo = ""
	}

	err = m.DB.UpdateAssistant(&assistant)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "/admin/assistant/" + strconv.Itoa(assistant.ID)
	case MANAGER_ACCESS_LEVEL:
		template = "/manager" + strconv.Itoa(assistant.ID)
	}

	m.App.Session.Put(r.Context(), "flash", "Ассистент успешно обновлен")

	http.Redirect(w, r, template, http.StatusSeeOther)
}
