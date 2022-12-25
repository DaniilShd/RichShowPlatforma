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

	http.Redirect(w, r, "/admin/store-all", http.StatusSeeOther)
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

func (m *Repository) AllStoreOrder(w http.ResponseWriter, r *http.Request) {
	src := chi.URLParam(r, "src")

	stringMap := make(map[string]string, 1)
	stringMap["type"] = src

	var storeOrders []models.StoreLead
	var err error

	switch src {
	case "new":
		storeOrders, err = m.DB.GetAllNewStoreOrder()
	case "completed":
		storeOrders, err = m.DB.GetAllCompleteStoreOrder()
	case "destroy":
		storeOrders, err = m.DB.GetAllToDestroyStoreOrder()
	}

	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["store-lead"] = storeOrders

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-all-store-leads.page.html"
	case STORE_ACCESS_LEVEL:
		template = "store.page.html"
	}

	render.Template(w, r, template, &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) ShowStoreOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := chi.URLParam(r, "src")

	stringMap := make(map[string]string, 1)
	stringMap["type"] = src

	storeOrder, err := m.DB.GetStoreOrderByID(id)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	checkList, err := m.DB.GetCheckListByID(storeOrder.CheckListID)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var items []models.Item

	for _, item := range checkList.Items {
		item.AmountItemOnce = float64(storeOrder.AmountOfChilds) * item.AmountItemOnce
		items = append(items, item)
	}

	checkList.Items = items

	data := make(map[string]interface{})
	data["store-lead"] = storeOrder
	data["check-list"] = checkList

	var template string

	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-order-show-new.page.html"
	case STORE_ACCESS_LEVEL:
		template = "store-order-show-new.page.html"
	}

	render.Template(w, r, template, &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
		Form:      forms.New(nil),
	})
}

func (m *Repository) ShowPostStoreOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := chi.URLParam(r, "src")

	err = r.ParseForm() // grab the multipart form
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	storeOrder, err := m.DB.GetStoreOrderByID(id)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	formdata := r.MultipartForm
	storeOrder.StoreDescription = r.Form.Get("description")
	storeOrder.ID = id

	file := formdata.File["photo"]

	if len(file) == 0 {
		form.Errors.Add("photo", "Добавьте фото!")
	} else {
		photo := file[0].Filename
		storeOrder.Photo = photo
	}
	if !form.Valid() {

		stringMap := make(map[string]string, 1)
		stringMap["type"] = src

		text := r.Form.Get("description")

		// storeOrder, err := m.DB.GetStoreOrderByID(id)
		// if err != nil {
		// 	fmt.Println(err)
		// 	helpers.ServerError(w, err)
		// 	return
		// }
		storeOrder.StoreDescription = text

		checkList, err := m.DB.GetCheckListByID(storeOrder.CheckListID)
		if err != nil {
			fmt.Println(err)
			helpers.ServerError(w, err)
			return
		}

		var items []models.Item

		for _, item := range checkList.Items {
			item.AmountItemOnce = float64(storeOrder.AmountOfChilds) * item.AmountItemOnce
			items = append(items, item)
		}

		checkList.Items = items

		data := make(map[string]interface{})
		data["store-lead"] = storeOrder
		data["check-list"] = checkList

		var template string
		switch src {
		case "new":
			switch m.App.Session.Get(r.Context(), "access_level") {
			case ADMIN_ACCESS_LEVEL:
				template = "admin-order-show-new.page.html"
			case STORE_ACCESS_LEVEL:
				template = "store-order-show-new.page.html"
			}
		case "completed":
		case "destroy":
		}

		render.Template(w, r, template, &models.TemplateData{
			Data:      data,
			StringMap: stringMap,
			Form:      form,
		})
		return
	}
	uuidPhoto := uuid.New().String()
	suffix := strings.SplitAfter(file[0].Filename, ".")[1]

	out, err := os.Create("./static/img/store-leads/" + uuidPhoto + "." + suffix)
	storeOrder.Photo = uuidPhoto + "." + suffix
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

	fmt.Println(storeOrder)

	err = m.DB.InsertStoreOrder(storeOrder)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "/admin/store-leads/" + src
	case STORE_ACCESS_LEVEL:
		template = "/" + src
	}

	m.App.Session.Put(r.Context(), "flash", "Заказ собран")

	http.Redirect(w, r, template, http.StatusSeeOther)
}

func (m *Repository) DestroyStoreOrder(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := chi.URLParam(r, "src")

	err = m.DB.DeleteStoreOrderByID(id)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "/admin/store-leads/" + src
	case STORE_ACCESS_LEVEL:
		template = "/" + src
	}

	m.App.Session.Put(r.Context(), "flash", "Заказ разобран")

	http.Redirect(w, r, template, http.StatusSeeOther)
}

func (m *Repository) ChangePostStoreOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := chi.URLParam(r, "src")

	err = r.ParseForm() // grab the multipart form
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	storeOrder, err := m.DB.GetStoreOrderByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	formdata := r.MultipartForm
	if r.Form.Get("description") != "" {
		storeOrder.StoreDescription = r.Form.Get("description")
	}

	file := formdata.File["photo"]

	if len(file) != 0 {

		if storeOrder.Photo != "" {
			filePath := "static/img/store-leads/" + storeOrder.Photo
			fmt.Println(filePath)
			err := os.Remove(filePath)
			if err != nil {
				helpers.ServerError(w, err)
				return
			}
		}

		uuidPhoto := uuid.New().String()
		suffix := strings.SplitAfter(file[0].Filename, ".")[1]

		out, err := os.Create("./static/img/store-leads/" + uuidPhoto + "." + suffix)
		storeOrder.Photo = uuidPhoto + "." + suffix
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

	}

	fmt.Println(storeOrder)

	err = m.DB.UpdateStoreOrder(storeOrder)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "/admin/store-lead/" + src + "/" + strconv.Itoa(id)
	case STORE_ACCESS_LEVEL:
		template = "/" + src
	}

	m.App.Session.Put(r.Context(), "flash", "Заказ изменен")

	http.Redirect(w, r, template, http.StatusSeeOther)
}
