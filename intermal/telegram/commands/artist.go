package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	modelsFromApp "github.com/DaniilShd/RichShowPlatforma/intermal/models"

	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *Commander) ArtistLeadsList(inputMesage *tgbotapi.Message) {

	leads, err := c.DB.GetAllLeadsOfArtistByChatID(inputMesage.Chat.ID)
	if err != nil {

		log.Fatal(err)
	}

	if len(leads) > 0 {
		keyBoard := tgbotapi.NewInlineKeyboardMarkup().InlineKeyboard

		count := 0
		number := 1
		var listButtonArray []tgbotapi.InlineKeyboardButton
		var res string

		for i, lead := range leads {
			data, _ := json.Marshal(models.RequestFromChat{
				Command: "get_lead",
				ChatID:  inputMesage.Chat.ID,
				LeadID:  lead.ID,
			})

			heroes := strings.Join(lead.NameHeroes, ", ")

			res = res + fmt.Sprintf("%d) № - %d, Дата: <strong>%s</strong>, Время: <strong>%s</strong>, герои: %s",
				number,
				lead.ID,
				string(lead.Date.Format("02-01-2006")),
				string(lead.Time.Format("15:04")),
				heroes) + "\n"

			if count != 3 {
				listButton := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(number), string(data))
				listButtonArray = append(listButtonArray, listButton)
				count++
			} else {
				keyBoard = append(keyBoard, listButtonArray)
				listButtonArray = nil
				count = 1
				listButton := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(number), string(data))
				listButtonArray = append(listButtonArray, listButton)
			}
			if i == len(leads)-1 {
				keyBoard = append(keyBoard, listButtonArray)
			}
			number++
		}
		msg := tgbotapi.NewMessage(inputMesage.Chat.ID, res)
		msg.ParseMode = "html"

		msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: keyBoard,
		}
		c.bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(inputMesage.Chat.ID, "Заказов нет")
		c.bot.Send(msg)
	}

}

func (c *Commander) ArtistLeadsToday(inputMesage *tgbotapi.Message) {

	leads, err := c.DB.GetAllLeadsOfArtistTodayByChatID(inputMesage.Chat.ID)
	if err != nil {

		log.Fatal(err)
	}

	if len(leads) > 0 {
		keyBoard := tgbotapi.NewInlineKeyboardMarkup().InlineKeyboard

		count := 0
		number := 1
		var listButtonArray []tgbotapi.InlineKeyboardButton
		var res string

		for i, lead := range leads {
			data, _ := json.Marshal(models.RequestFromChat{
				Command: "get_lead",
				ChatID:  inputMesage.Chat.ID,
				LeadID:  lead.ID,
			})

			heroes := strings.Join(lead.NameHeroes, ", ")

			res = res + fmt.Sprintf("%d) № - %d, Дата: <strong>%s</strong>, Время: <strong>%s</strong>, герои: %s",
				number,
				lead.ID,
				string(lead.Date.Format("02-01-2006")),
				string(lead.Time.Format("15:04")),
				heroes) + "\n"

			if count != 3 {
				listButton := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(number), string(data))
				listButtonArray = append(listButtonArray, listButton)
				count++
			} else {
				keyBoard = append(keyBoard, listButtonArray)
				listButtonArray = nil
				count = 1
				listButton := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(number), string(data))
				listButtonArray = append(listButtonArray, listButton)
			}
			if i == len(leads)-1 {
				keyBoard = append(keyBoard, listButtonArray)
			}
			number++
		}
		msg := tgbotapi.NewMessage(inputMesage.Chat.ID, res)
		msg.ParseMode = "html"

		msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: keyBoard,
		}
		c.bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(inputMesage.Chat.ID, "Заказов нет")
		c.bot.Send(msg)
	}

}

func (c *Commander) GetLeadByID(request models.RequestFromChat) {
	var res string
	var checkLead string
	var masterClass string
	var animation string
	var partyAndQuest string
	var show string
	var other string
	var heroes string
	var assistant string
	var gender string

	ResponseFromApp := make(chan interface{})
	request.ResponseLeadFromApp = ResponseFromApp

	c.App.RequestFromTelegram <- request

	lead := <-request.ResponseLeadFromApp

	leadMy := lead.(modelsFromApp.Lead)

	clientPhoneNumber := strings.ReplaceAll(leadMy.Client.PhoneNumber, " ", "")
	clientPhoneNumber = strings.ReplaceAll(clientPhoneNumber, ")", "")
	clientPhoneNumber = strings.ReplaceAll(clientPhoneNumber, "-", "")
	clientPhoneNumber = strings.ReplaceAll(clientPhoneNumber, "(", "")
	clientPhoneNumber = "<a href=\"tel:" + clientPhoneNumber + "\">" + clientPhoneNumber + "</a>"

	if leadMy.Confirmed {
		checkLead = "Да"
	} else {
		checkLead = "Нет"
	}

	if leadMy.Child.Gender == 1 {
		gender = "мужской"
	} else {
		gender = "женский"
	}

	for _, item := range leadMy.MasterClasses {
		masterClass = masterClass + item.Name + " " + "(" + item.Description + ")" + ", "
	}

	for _, item := range leadMy.Shows {
		show = show + item.Name + " " + "(" + item.Description + ")" + ", "
	}

	for _, item := range leadMy.PartyAndQuests {
		partyAndQuest = partyAndQuest + item.Name + " " + "(" + item.Description + ")" + ", "
	}

	for _, item := range leadMy.Animations {
		animation = animation + item.Name + " " + "(" + item.Description + ")" + ", "
	}

	for _, item := range leadMy.Others {
		other = other + item.Name + " " + "(" + item.Description + ")" + ", "
	}

	for _, item := range leadMy.Assistants {
		assistant = assistant + item.FirstName + " " + item.LastName + " (Телефон: <a href=\"tel:+7" + item.PhoneNumber + "\">+7" + item.PhoneNumber + "</a>)" + ", "
	}
	var heroesRes string
	for _, item := range leadMy.Heroes {
		hero := item.HeroName + " " + "(Аниматор: " + item.ArtistFirstName + " " + item.ArtistLastName + ", Телефон: <a href=\"tel:+7" + item.PhoneNumber + "\">+7" + item.PhoneNumber + "</a>)"
		heroesRes = strings.Join([]string{heroes, hero}, ", ")
	}

	res = fmt.Sprintf(`
	&#10024; <strong>№Заказа:</strong> %d
	&#128198; <strong>Дата: %s</strong>
	&#128348; <strong>Время: %s</strong>
	&#128313; <strong>Адрес:</strong> %s
	&#128313; <strong>Герои:</strong> %s
	&#128313; <strong>Имя клиента:</strong> %s
	&#128313; <strong>Имя ребенка:</strong> %s
	&#128313; <strong>Возраст ребенка:</strong> %d
	&#128313; <strong>Пол ребенка:</strong> %s
	&#128222; <strong>Номер:</strong> %s
	&#128313; <strong>Продолжительность:</strong> %d мин
	&#128313; <strong>Среднее количество детей:</strong> %d
	&#128313; <strong>Средний возраст детей:</strong> %d лет
	&#128313; <strong>Заказ подтвержден:</strong> %s
	&#128313; <strong>Комментарий:</strong> %s
	&#128313; <strong>Мастер-классы:</strong> %s
	&#128313; <strong>Анимация:</strong> %s
	&#128313; <strong>Шоу-программы:</strong> %s
	&#128313; <strong>Вечеринки и квесты:</strong> %s
	&#128313; <strong>Другое:</strong> %s
	&#128313; <strong>Ассистенты:</strong> %s
	`, leadMy.ID,
		string(leadMy.Date.Format("02-01-2006")),
		string(leadMy.Time.Format("15:04")),
		leadMy.Address,
		heroesRes,
		fmt.Sprintf("%s %s", leadMy.Client.FirstName, leadMy.Client.LastName),
		leadMy.Child.Name,
		leadMy.Child.Age,
		gender,
		clientPhoneNumber,
		leadMy.Duration,
		leadMy.AmountOfChildren,
		leadMy.AverageAgeOfChildren,
		checkLead,
		leadMy.Description,
		masterClass,
		animation,
		show,
		partyAndQuest,
		other,
		assistant)

	data, _ := json.Marshal(models.RequestFromChat{
		Command: "get_order_by_id",
		ChatID:  request.ChatID,
		LeadID:  leadMy.ID,
	})

	msg := tgbotapi.NewMessage(request.ChatID, res)
	msg.ParseMode = "html"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Реквизит", string(data)),
		),
	)
	c.bot.Send(msg)
}

func (c *Commander) GetOrderByID(request models.RequestFromChat) {
	idOrderStore, err := c.DB.GetOrderStoreIDByLeadID(request.LeadID)
	if err != nil {
		log.Fatal(err)
	}

	ResponseFromApp := make(chan interface{})
	request.ResponseLeadFromApp = ResponseFromApp

	request.StoreOrderID = idOrderStore

	c.App.RequestFromTelegram <- request

	orderStore := <-request.ResponseLeadFromApp

	res := orderStore.([]modelsFromApp.StoreLead)

	var CheckListRequest models.RequestFromChat
	CheckListRequest.Command = "get_checklist"

	for _, storeOrder := range res {
		CheckListRequest.CheckListID = storeOrder.CheckListID

		if storeOrder.Completed && !storeOrder.Canceled {

			ResponseFromApp := make(chan interface{})
			CheckListRequest.ResponseLeadFromApp = ResponseFromApp

			c.App.RequestFromTelegram <- CheckListRequest
			checkList := <-CheckListRequest.ResponseLeadFromApp
			checkListRes := checkList.(modelsFromApp.CheckList)
			file := tgbotapi.FilePath("./static/img/store-leads/" + storeOrder.Photo)
			res := tgbotapi.NewPhoto(request.ChatID, file)

			text := "<strong>Название программы</strong> - " + checkListRes.Name + "\n"
			text = text + "<strong>Список реквизита</strong>" + "\n"

			for _, item := range checkListRes.NameOfPoints {
				text = text + "&#9642; " + item + "\n"
			}
			text = text + "<strong>Список расходников</strong>" + "\n"

			for _, item := range checkListRes.Items {
				text = text + "&#9642; " + item.Name + " - " + strconv.Itoa(int(item.AmountItemOnce*float64(storeOrder.AmountOfChilds))) + " " + item.Dimension + "\n"
			}

			text = text + "<strong>Комментарий:</strong> " + storeOrder.StoreDescription

			res.ParseMode = "html"
			fmt.Println()
			res.Caption = text
			c.bot.Send(res)
		} else if !storeOrder.Canceled {

			text := "Название программы: " + storeOrder.Name + "\n" + "Реквизит еще не собран"
			msg := tgbotapi.NewMessage(request.ChatID, text)
			c.bot.Send(msg)
		}

	}

}
