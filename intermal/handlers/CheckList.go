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

// Check lists all
func (m *Repository) CheckListAll(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")
	src := exploded[3]
	var typeOfCheckList int
	var title string
	switch src {
	case "program-show":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_SHOW
		title = "Шоу программы"
	case "class-master":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_MASTER_CLASS
		title = "Мастер-классы"
	case "animcheck":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_ANIMATION
		title = "Анимационные программы"
	case "parties":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_PARTIES_AND_QUESTS
		title = "Вечеринки и квесты"
	case "oths":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_OTHER
		title = "Другое"
	}

	checkLists, err := m.DB.GetAllCheckListsOfType(typeOfCheckList)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["check-lists"] = checkLists
	stringMap := make(map[string]string)
	stringMap["source"] = src
	stringMap["title"] = title

	render.Template(w, r, "admin-all-check-lists.page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// Create new ckeck-list
func (m *Repository) NewCheckList(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")
	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["source"] = src

	checkList := models.CheckList{}

	var typeOfCheckList int
	switch src {
	case "program-show":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_SHOW
	case "class-master":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_MASTER_CLASS
	case "animcheck":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_ANIMATION
	case "parties":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_PARTIES_AND_QUESTS
	case "oths":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_OTHER
	}

	var storeItems []models.StoreItem
	storeItems, err := m.DB.GetAllStoreItem()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["check-list"] = checkList
	data["type-of-list"] = typeOfCheckList
	data["store-items"] = storeItems

	render.Template(w, r, "admin-check-list-new.page.html", &models.TemplateData{
		Data:      data,
		Form:      forms.New(nil),
		StringMap: stringMap,
	})
}

func (m *Repository) NewPostCheckList(w http.ResponseWriter, r *http.Request) {

	//Достаем значение из строки URL
	exploded := strings.Split(r.RequestURI, "/")
	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["source"] = src

	fmt.Println(src)

	//Парсим форму
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant't parse form!")
		http.Redirect(w, r, "admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name_check_list", "description")

	//Сохраняем данные в сущности checkList
	var checkList models.CheckList
	checkList.Duration, _ = strconv.Atoi(r.Form.Get("duration"))
	checkList.Name = r.Form.Get("name_check_list")
	checkList.Description = r.Form.Get("description")
	checkList.TypeOfList, _ = strconv.Atoi(r.Form.Get("id_type_of_list"))
	NameOfPoints := r.Form["name_of_points[]"]
	checkList.NameOfPoints = append(checkList.NameOfPoints, NameOfPoints...)

	//Store item здесь парсим данные которые касаются выбранных материалов
	itemsID := r.Form["check_list_store[]"]
	itemsMinAmount := r.Form["amount_item_once[]"]

	//Здесь сохраняем все выбранные материалы, а также проверяем правильно ли были введены числа
	var items []models.Item
	for i := 1; i < len(itemsID); i++ {
		var item models.Item
		id, err := strconv.Atoi(itemsID[i])
		if err != nil {
			continue
		}
		item.ID = id
		storeItem, _ := m.DB.GetStoreItemByID(item.ID)
		item.Dimension = storeItem.Dimension
		item.Name = storeItem.Name
		item.AmountItemOnce, err = strconv.ParseFloat(itemsMinAmount[i], 64)
		if err != nil {
			item.AmountItemOnce = 0 //itemsMinAmount[i]
			form.Errors.Add("amount_item_once", "Необходимо ввести число, десятичные числа вводить через точку!")
		}
		items = append(items, item)
	}
	checkList.Items = items

	if !form.Valid() {
		var storeItems []models.StoreItem
		storeItems, err := m.DB.GetAllStoreItem()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		data := make(map[string]interface{})
		data["check-list"] = checkList
		data["type-of-list"] = checkList.TypeOfList
		data["store-items"] = storeItems

		render.Template(w, r, "admin-check-list-new.page.html", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		})
		return
	}

	err = m.DB.InsertCheckList(&checkList)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Чек-лист успешно сохранен")

	http.Redirect(w, r, fmt.Sprintf("/admin/check-lists/%s", src), http.StatusSeeOther)
}

//Delete master-class Chek-List

func (m *Repository) DeleteCheсkList(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")
	src := exploded[3]

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.DB.DeleteCheckListByID(id)
	if err != nil {
		fmt.Println(err)
		m.App.Session.Put(r.Context(), "error", "Чек лист используется")
		http.Redirect(w, r, fmt.Sprintf("/admin/check-lists/%s", src), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/check-lists/%s", src), http.StatusSeeOther)
}

func (m *Repository) ShowCheckList(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["source"] = src

	var typeOfCheckList int
	switch src {
	case "show-program":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_SHOW
	case "class-master":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_MASTER_CLASS
	case "animcheck":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_ANIMATION
	case "parties":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_PARTIES_AND_QUESTS
	case "oths":
		typeOfCheckList = CHECK_LISTS_TYPE_OF_OTHER
	}

	var storeItems []models.StoreItem
	storeItems, err = m.DB.GetAllStoreItem()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	checkList, err := m.DB.GetCheckListByID(id)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["check-list"] = checkList
	data["type-of-list"] = typeOfCheckList
	data["store-items"] = storeItems

	render.Template(w, r, "admin-check-list-show.page.html", &models.TemplateData{
		Data:      data,
		Form:      forms.New(nil),
		StringMap: stringMap,
	})
}

func (m *Repository) ShowPostCheckList(w http.ResponseWriter, r *http.Request) {

	//Достаем значение из строки URL
	exploded := strings.Split(r.RequestURI, "/")
	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["source"] = src

	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//Парсим форму
	err = r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant't parse form!")
		http.Redirect(w, r, "admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	checkList, err := m.DB.GetCheckListByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name_check_list", "description")

	//Сохраняем данные в сущности checkList
	// var checkList models.CheckList
	checkList.Name = r.Form.Get("name_check_list")
	checkList.Duration, _ = strconv.Atoi(r.Form.Get("duration"))
	checkList.Description = r.Form.Get("description")
	checkList.TypeOfList, _ = strconv.Atoi(r.Form.Get("id_type_of_list"))
	checkList.NameOfPoints = nil
	NameOfPoints := r.Form["name_of_points[]"]
	checkList.NameOfPoints = append(checkList.NameOfPoints, NameOfPoints...)

	//Store item здесь парсим данные которые касаются выбранных материалов
	itemsID := r.Form["check_list_store[]"]
	itemsMinAmount := r.Form["amount_item_once[]"]

	//Здесь сохраняем все выбранные материалы, а также проверяем правильно ли были введены числа
	var items []models.Item
	for i := 1; i < len(itemsID); i++ {
		var item models.Item
		id, err := strconv.Atoi(itemsID[i])
		//Пропускаем материал который не был указан в форме
		if err != nil {
			continue
		}
		item.ID = id
		storeItem, _ := m.DB.GetStoreItemByID(item.ID)
		item.Dimension = storeItem.Dimension
		item.Name = storeItem.Name
		item.AmountItemOnce, err = strconv.ParseFloat(itemsMinAmount[i], 64)
		if err != nil {
			item.AmountItemOnce = 0 //itemsMinAmount[i]
			form.Errors.Add("amount_item_once", "Необходимо ввести число, десятичные числа вводить через точку!")
		}
		items = append(items, item)
	}
	checkList.Items = items

	if !form.Valid() {

		var storeItems []models.StoreItem
		storeItems, err = m.DB.GetAllStoreItem()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		data := make(map[string]interface{})
		data["check-list"] = checkList
		data["store-items"] = storeItems

		render.Template(w, r, "admin-check-list-show.page.html", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		})
		return
	}

	err = m.DB.UpdateCheckList(checkList)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Чек-лист успешно сохранен")

	http.Redirect(w, r, fmt.Sprintf("/admin/check-lists/%s", src), http.StatusSeeOther)
}
