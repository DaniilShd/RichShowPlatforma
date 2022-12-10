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

const (
	ADMIN_ACCESS_LEVEL   = 1
	MANAGER_ACCESS_LEVEL = 2
	STORE_ACCESS_LEVEL   = 3
)

func (m *Repository) StoreItemAll(w http.ResponseWriter, r *http.Request) {
	storeItem, err := m.DB.GetAllStoreItem()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["store"] = storeItem

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-all-store.page.html"
	case STORE_ACCESS_LEVEL:
		template = "store.page.html"
	}

	render.Template(w, r, template, &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) NewStoreItem(w http.ResponseWriter, r *http.Request) {
	storeItem := models.StoreItem{}
	data := make(map[string]interface{})
	data["store"] = storeItem

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-store-new.page.html"
	case STORE_ACCESS_LEVEL:
		template = "store-new.page.html"
	}

	render.Template(w, r, template, &models.TemplateData{
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

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-store-new.page.html"
	case STORE_ACCESS_LEVEL:
		template = "store-new.page.html"
	}

	if !form.Valid() {

		data := make(map[string]interface{})
		data["store"] = storeItem

		render.Template(w, r, template, &models.TemplateData{
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

	m.App.Session.Put(r.Context(), "flash", "Материал сохранен")

	http.Redirect(w, r, template, http.StatusSeeOther)
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

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-store-show.page.html"
	case STORE_ACCESS_LEVEL:
		template = "store-show.page.html"
	}

	render.Template(w, r, template, &models.TemplateData{
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

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-store-show.page.html"
	case STORE_ACCESS_LEVEL:
		template = "store-show.page.html"
	}

	if !form.Valid() {

		data := make(map[string]interface{})
		data["store"] = storeItem

		render.Template(w, r, template, &models.TemplateData{
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

	m.App.Session.Put(r.Context(), "flash", "Материал изменен")

	http.Redirect(w, r, template, http.StatusSeeOther)
}

func (m *Repository) DeleteStoreItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete item")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "/admin/store"
	case STORE_ACCESS_LEVEL:
		template = "/store"
	}

	err = m.DB.DeleteStoreItemByID(id)
	if err != nil {
		fmt.Println(err)
		m.App.Session.Put(r.Context(), "error", "Не удалось удалить (материал используется в чек листе)")

		http.Redirect(w, r, template, http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "success", "Успешно удалено")
	http.Redirect(w, r, template, http.StatusSeeOther)
}

// var count = 0

// func (m *Repository) TestFetch(w http.ResponseWriter, r *http.Request) {
// 	count++

// 	js, _ := json.Marshal(count)
// 	fmt.Println(js)

// 	w.Write(js)
// }
