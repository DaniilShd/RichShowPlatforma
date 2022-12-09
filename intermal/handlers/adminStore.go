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

func (m *Repository) StoreItemAll(w http.ResponseWriter, r *http.Request) {
	storeItem, err := m.DB.GetAllStoreItem()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["store"] = storeItem

	render.Template(w, r, "admin-all-store.page.html", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) NewStoreItem(w http.ResponseWriter, r *http.Request) {
	storeItem := models.StoreItem{}
	data := make(map[string]interface{})
	data["store"] = storeItem

	render.Template(w, r, "admin-store-new.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) NewPostStoreItem(w http.ResponseWriter, r *http.Request) {
	var storeItem models.StoreItem
	storeItem.Name = r.Form.Get("name_item")
	storeItem.CurrentAmount, _ = strconv.Atoi(r.Form.Get("current_amount"))
	storeItem.MinAmount, _ = strconv.Atoi(r.Form.Get("min_amount"))
	storeItem.Description = r.Form.Get("description")
	storeItem.Dimension = r.Form.Get("dimension")

	form := forms.New(r.PostForm)
	form.Required("name_item", "current_amount", "min_amount", "dimension")

	if !form.Valid() {

		data := make(map[string]interface{})
		data["store"] = storeItem

		render.Template(w, r, "admin-store-new.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err := m.DB.InsertStoreItem(&storeItem)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Store saved")

	http.Redirect(w, r, "/admin/store", http.StatusSeeOther)
}

func (m *Repository) ShowStoreItem(w http.ResponseWriter, r *http.Request) {

	//Достаем значение из строки URL
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[3])
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	storeItem, err := m.DB.GetStoreItemByID(id)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["store"] = storeItem

	render.Template(w, r, "admin-store-show.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) ShowPostStoreItem(w http.ResponseWriter, r *http.Request) {

	//Достаем значение из строки URL
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[3])
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var storeItem models.StoreItem
	storeItem.ID = id
	storeItem.Name = r.Form.Get("name_item")
	storeItem.CurrentAmount, _ = strconv.Atoi(r.Form.Get("current_amount"))
	storeItem.MinAmount, _ = strconv.Atoi(r.Form.Get("min_amount"))
	storeItem.Description = r.Form.Get("description")
	storeItem.Dimension = r.Form.Get("dimension")

	form := forms.New(r.PostForm)
	form.Required("name_item", "current_amount", "min_amount", "dimension")

	if !form.Valid() {

		data := make(map[string]interface{})
		data["store"] = storeItem

		render.Template(w, r, "admin-store-show.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.UpdateStoreItem(&storeItem)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Store changed")

	http.Redirect(w, r, "/admin/store", http.StatusSeeOther)
}

func (m *Repository) DeleteStoreItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.DB.DeleteStoreItemByID(id)
	if err != nil {
		fmt.Println(err)
		m.App.Session.Put(r.Context(), "error", "Не удалось удалить")
		http.Redirect(w, r, "/admin/store", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "success", "Успешно удалено")
	http.Redirect(w, r, "/admin/store", http.StatusSeeOther)
}
