package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/forms"
	"github.com/DaniilShd/RichShowPlatforma/intermal/helpers"
	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
	"github.com/DaniilShd/RichShowPlatforma/intermal/render"
	"github.com/go-chi/chi"
)

func (m *Repository) NewLead(w http.ResponseWriter, r *http.Request) {

	lead := models.Lead{}
	data := make(map[string]interface{})
	lead.Date = time.Now()
	lead.Child.DateOfBirthDay = time.Now()
	//Выбираем из базы данных все программы
	shows, err := m.DB.GetAllCheckListsOfType(CHECK_LISTS_TYPE_OF_SHOW)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	master_class, err := m.DB.GetAllCheckListsOfType(CHECK_LISTS_TYPE_OF_MASTER_CLASS)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	animation, err := m.DB.GetAllCheckListsOfType(CHECK_LISTS_TYPE_OF_ANIMATION)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	party_and_quest, err := m.DB.GetAllCheckListsOfType(CHECK_LISTS_TYPE_OF_PARTIES_AND_QUESTS)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	other, err := m.DB.GetAllCheckListsOfType(CHECK_LISTS_TYPE_OF_OTHER)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	//Берем из базы данных всех артистов
	artists, err := m.DB.GetAllArtists()
	if err != nil {
		fmt.Println(err, "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		helpers.ServerError(w, err)
		return
	}
	//Берем из базы данных всех ассистентов
	assistants, err := m.DB.GetAllAssistants()
	if err != nil {
		fmt.Println(err, "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		helpers.ServerError(w, err)
		return
	}
	//Берем из базы данных всех героев
	heroes, err := m.DB.GetAllHeroes()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	//передаем все в шаблон
	data["lead"] = lead
	data["shows"] = shows
	data["master_class"] = master_class
	data["animation"] = animation
	data["party_and_quest"] = party_and_quest
	data["other"] = other
	data["artists"] = artists
	data["assistants"] = assistants
	data["heroes"] = heroes

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-lead-new.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "manager-new.page.html"
	}

	render.Template(w, r, template, &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) NewPostLead(w http.ResponseWriter, r *http.Request) {

	//Парсим форму
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant't parse form!")
		http.Redirect(w, r, "admin/leads", http.StatusTemporaryRedirect)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "phone_number", "name_child", "time", "address")
	form.MinLength("phone_number", 17, r)
	//Сохраняем данные в сущности lead
	var lead models.Lead
	lead.Client.FirstName = r.Form.Get("first_name")
	lead.Client.LastName = r.Form.Get("last_name")
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
	lead.Client.PhoneNumber = phoneNumber

	lead.Client.Telegram = r.Form.Get("telegram_client")
	lead.Child.Name = r.Form.Get("name_child")
	lead.Child.Age, err = strconv.Atoi(r.Form.Get("age"))
	if err != nil {
		form.Errors.Add("age", "Необходимо ввести число!")
	}

	dateForm := r.Form.Get("date")

	layoutDateForm := "02-01-2006"

	dateRes, err := time.Parse(layoutDateForm, dateForm)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	timeForm := r.Form.Get("time")
	if timeForm == "00:00" {
		form.Errors.Add("time", "Введите время")
	}
	layoutTimeForm := "15:04"
	timeRes, err := time.Parse(layoutTimeForm, timeForm)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	lead.Date = dateRes
	lead.Time = timeRes
	lead.Address = r.Form.Get("address")
	dateBirth := r.Form.Get("date_of_birthday_child")
	dateLayout := "02-01-2006"
	dateBirthRes, err := time.Parse(dateLayout, dateBirth)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}
	lead.Child.DateOfBirthDay = dateBirthRes
	lead.AmountOfChildren, err = strconv.Atoi(r.Form.Get("amount_of_children"))
	if err != nil {
		form.Errors.Add("amount_of_children", "Необходимо ввести число!")
		err = nil
	}
	lead.AverageAgeOfChildren, err = strconv.Atoi(r.Form.Get("average_age_of_children"))
	if err != nil {
		form.Errors.Add("average_age_of_children", "Необходимо ввести число!")
		err = nil
	}

	lead.Child.Gender, _ = strconv.Atoi(r.Form.Get("gender_child"))

	lead.Description = r.Form.Get("description")

	//Выбираем ids шоу классы указанные в форме вместе с комментарием---------------------------
	showIDs := r.Form["shows[]"]
	showsDescription := r.Form["shows_description[]"]

	//Здесь сохраняем все выбранные шоу программы
	var shows []models.Show
	for i := 1; i < len(showIDs); i++ {
		var show models.Show
		id, err := strconv.Atoi(showIDs[i])
		if err != nil {
			continue
		}
		show.ID = id
		checkListItem, _ := m.DB.GetCheckListByID(show.ID)
		show.Duration = checkListItem.Duration
		show.Name = checkListItem.Name
		show.Description = showsDescription[i]
		shows = append(shows, show)
		lead.Duration = lead.Duration + checkListItem.Duration
	}
	lead.Shows = shows
	//---------------------------------------------
	masterClassIDs := r.Form["master_class[]"]
	masterClassDescription := r.Form["master_class_description[]"]

	//Здесь сохраняем все выбранные шоу программы
	var masterClasses []models.MasterClass
	for i := 1; i < len(masterClassIDs); i++ {
		var masterClass models.MasterClass
		id, err := strconv.Atoi(masterClassIDs[i])
		if err != nil {
			continue
		}
		masterClass.ID = id
		checkListItem, _ := m.DB.GetCheckListByID(masterClass.ID)
		masterClass.Duration = checkListItem.Duration
		masterClass.Name = checkListItem.Name
		masterClass.Description = masterClassDescription[i]
		masterClasses = append(masterClasses, masterClass)
		lead.Duration = lead.Duration + checkListItem.Duration
	}
	lead.MasterClasses = masterClasses
	//---------------------------------------------
	animationIDs := r.Form["animation[]"]
	animationDescription := r.Form["animation_description[]"]

	//Здесь сохраняем все выбранные шоу программы
	var animations []models.Animation
	for i := 1; i < len(animationIDs); i++ {
		var animation models.Animation
		id, err := strconv.Atoi(animationIDs[i])
		if err != nil {
			continue
		}
		animation.ID = id
		checkListItem, _ := m.DB.GetCheckListByID(animation.ID)
		animation.Duration = checkListItem.Duration
		animation.Name = checkListItem.Name
		animation.Description = animationDescription[i]
		animations = append(animations, animation)
		lead.Duration = lead.Duration + checkListItem.Duration
	}
	lead.Animations = animations

	//---------------------------------------------
	partyIDs := r.Form["party_and_quest[]"]
	partyDescription := r.Form["party_and_quest_description[]"]

	//Здесь сохраняем все выбранные шоу программы
	var partys []models.PartyAndQuest
	for i := 1; i < len(partyIDs); i++ {
		var party models.PartyAndQuest
		id, err := strconv.Atoi(partyIDs[i])
		if err != nil {
			continue
		}
		party.ID = id
		checkListItem, _ := m.DB.GetCheckListByID(party.ID)
		party.Duration = checkListItem.Duration
		party.Name = checkListItem.Name
		party.Description = partyDescription[i]
		partys = append(partys, party)
		lead.Duration = lead.Duration + checkListItem.Duration
	}
	lead.PartyAndQuests = partys

	//---------------------------------------------
	otherIDs := r.Form["other[]"]
	otherDescription := r.Form["other_description[]"]

	//Здесь сохраняем все выбранные шоу программы
	var others []models.Other
	for i := 1; i < len(otherIDs); i++ {
		var other models.Other
		id, err := strconv.Atoi(otherIDs[i])
		if err != nil {
			continue
		}
		other.ID = id
		checkListItem, _ := m.DB.GetCheckListByID(other.ID)
		other.Duration = checkListItem.Duration
		other.Name = checkListItem.Name
		other.Description = otherDescription[i]
		others = append(others, other)
		lead.Duration = lead.Duration + checkListItem.Duration
	}
	lead.Others = others

	//Сохранение в лиде информации по героям и аниматорам
	heroesIDs := r.Form["id_hero[]"]
	artistID := r.Form["id_artist[]"]
	var heroes []models.LeadHero
	for i := 1; i < len(heroesIDs); i++ {
		var hero models.LeadHero
		idHero, err := strconv.Atoi(heroesIDs[i])
		if err != nil {
			continue
		}
		idArtist, err := strconv.Atoi(artistID[i])
		if err != nil {
			hero.ArtistID = idArtist
		} else {
			hero.ArtistID = idArtist
			artistDB, err := m.DB.GetArtistByID(hero.ArtistID)
			if err != nil {
				helpers.ServerError(w, err)
				return
			}
			hero.ArtistFirstName = artistDB.FirstName
			hero.ArtistLastName = artistDB.LastName
			fmt.Println(artistDB)
		}

		hero.HeroID = idHero
		heroDB, _ := m.DB.GetHeroByID(hero.HeroID)
		fmt.Println(heroDB)
		hero.HeroName = heroDB.Name
		heroes = append(heroes, hero)
	}
	lead.Heroes = heroes

	//Сохранение информации по ассистентам
	assistantIDs := r.Form["id_assistant[]"]

	var assistants []models.Assistant
	for i := 1; i < len(assistantIDs); i++ {
		var assistant models.Assistant
		id, err := strconv.Atoi(assistantIDs[i])
		if err != nil {
			continue
		}
		assistant.ID = id
		assistantItem, _ := m.DB.GetAssistantByID(assistant.ID)
		assistant.FirstName = assistantItem.FirstName
		assistant.LastName = assistantItem.LastName
		assistants = append(assistants, assistant)
	}
	lead.Assistants = assistants

	//Здесь сохраняем все выбранные шоу программы
	if !form.Valid() {
		//Выбираем из базы данных все программы
		shows, err := m.DB.GetAllCheckListsOfType(CHECK_LISTS_TYPE_OF_SHOW)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		master_class, err := m.DB.GetAllCheckListsOfType(CHECK_LISTS_TYPE_OF_MASTER_CLASS)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		animation, err := m.DB.GetAllCheckListsOfType(CHECK_LISTS_TYPE_OF_ANIMATION)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		party_and_quest, err := m.DB.GetAllCheckListsOfType(CHECK_LISTS_TYPE_OF_PARTIES_AND_QUESTS)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		other, err := m.DB.GetAllCheckListsOfType(CHECK_LISTS_TYPE_OF_OTHER)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		//Берем из базы данных всех артистов
		artists, err := m.DB.GetAllArtists()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		//Берем из базы данных всех ассистентов
		assistants, err := m.DB.GetAllAssistants()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		//Берем из базы данных всех героев
		heroes, err := m.DB.GetAllHeroes()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		lead.Client.PhoneNumber = helpers.ConvertNumberPhone(lead.Client.PhoneNumber)

		fmt.Println(lead)

		data := make(map[string]interface{})
		data["lead"] = lead
		data["shows"] = shows
		data["master_class"] = master_class
		data["animation"] = animation
		data["party_and_quest"] = party_and_quest
		data["other"] = other
		data["artists"] = artists
		data["assistants"] = assistants
		data["heroes"] = heroes

		var template string
		switch m.App.Session.Get(r.Context(), "access_level") {
		case ADMIN_ACCESS_LEVEL:
			template = "admin-lead-new.page.html"
		case MANAGER_ACCESS_LEVEL:
			template = "admin-lead-new.page.html"
		}

		render.Template(w, r, template, &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	fmt.Println(lead)

	err = m.DB.InsertLead(&lead)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Лид успешно сохранен")

	var link string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		link = "/admin/leads"
	case MANAGER_ACCESS_LEVEL:
		link = "/admin/leads"
	}

	http.Redirect(w, r, link, http.StatusSeeOther)
}

func (m *Repository) AllConfirmedLead(w http.ResponseWriter, r *http.Request) {
	leads, err := m.DB.GetAllConfirmedLeads()
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	stringMap := make(map[string]string, 1)
	stringMap["add"] = "false"

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-all-leads.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "all-leads.page.html"
	}

	data := make(map[string]interface{})
	data["leads"] = leads

	render.Template(w, r, template, &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) AllRawLead(w http.ResponseWriter, r *http.Request) {
	leads, err := m.DB.GetAllRawLeads()
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	stringMap := make(map[string]string, 1)
	stringMap["add"] = "true"

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-all-leads.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "all-leads.page.html"
	}

	data := make(map[string]interface{})
	data["leads"] = leads

	render.Template(w, r, template, &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) AllArchiveLead(w http.ResponseWriter, r *http.Request) {
	leads, err := m.DB.GetAllArchiveLeads()
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	stringMap := make(map[string]string, 1)
	stringMap["add"] = "false"

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-all-leads.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "all-leads.page.html"
	}

	data := make(map[string]interface{})
	data["leads"] = leads

	render.Template(w, r, template, &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) ShowLead(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	lead, err := m.DB.GetLeadByID(id)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["lead"] = lead

	var template string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		template = "admin-lead-show.page.html"
	case MANAGER_ACCESS_LEVEL:
		template = "lead-show.page.html"
	}

	render.Template(w, r, template, &models.TemplateData{
		Data: data,
	})
}

// Set and delete confirmed

func (m *Repository) SetConfirmedLead(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	err = m.DB.SetConfirmedLeadByID(id)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var link string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		link = "/admin/show-lead/" + strconv.Itoa(id)
	case MANAGER_ACCESS_LEVEL:
		link = "/manager/show-lead/" + strconv.Itoa(id)
	}

	m.App.Session.Put(r.Context(), "flash", "Заказ подтвержден")

	http.Redirect(w, r, link, http.StatusSeeOther)

}
func (m *Repository) DeleteConfirmedLead(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	err = m.DB.DeleteConfirmedLeadByID(id)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var link string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		link = "/admin/show-lead/" + strconv.Itoa(id)
	case MANAGER_ACCESS_LEVEL:
		link = "/manager/show-lead/" + strconv.Itoa(id)
	}

	m.App.Session.Put(r.Context(), "flash", "Подтверждение снято")

	http.Redirect(w, r, link, http.StatusSeeOther)

}

func (m *Repository) DeleteLead(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	err = m.DB.DeleteLeadByID(id)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	var link string
	switch m.App.Session.Get(r.Context(), "access_level") {
	case ADMIN_ACCESS_LEVEL:
		link = "/admin/rawleads"
	case MANAGER_ACCESS_LEVEL:
		link = "/manager/rawleads"
	}

	m.App.Session.Put(r.Context(), "flash", "Заказ удален")

	http.Redirect(w, r, link, http.StatusSeeOther)
}
