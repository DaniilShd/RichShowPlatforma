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
				listButton := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(lead.ID), string(data))
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
				listButton := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(lead.ID), string(data))
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

	ResponseFromApp := make(chan *modelsFromApp.Lead)
	request.ResponseLeadFromApp = ResponseFromApp

	c.App.RequestFromTelegram <- request

	leadMy := <-request.ResponseLeadFromApp

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

	for _, item := range leadMy.MasterClasses {
		masterClass = masterClass + item.Name + " " + "(Комментарий:" + item.Description + ")" + ", "
	}

	for _, item := range leadMy.Shows {
		show = show + item.Name + " " + "(Комментарий:" + item.Description + ")" + ", "
	}

	for _, item := range leadMy.PartyAndQuests {
		partyAndQuest = partyAndQuest + item.Name + " " + "(Комментарий:" + item.Description + ")" + ", "
	}

	for _, item := range leadMy.Animations {
		animation = animation + item.Name + " " + "(Комментарий:" + item.Description + ")" + ", "
	}

	for _, item := range leadMy.Others {
		other = other + item.Name + " " + "(Комментарий: " + item.Description + ")" + ", "
	}

	for _, item := range leadMy.Assistants {
		assistant = assistant + item.FirstName + " " + item.LastName + " (Телефон: <a href=\"tel:+7" + item.PhoneNumber + "\">+7" + item.PhoneNumber + "</a>)" + ", "
	}

	for _, item := range leadMy.Heroes {
		heroes = heroes + item.HeroName + " " + "(Аниматор: " + item.ArtistFirstName + " " + item.ArtistLastName + ")" + ", "
	}

	res = fmt.Sprintf(`
	&#10024; <strong>№Заказа:</strong> %d
	&#128198; <strong>Дата: %s</strong>
	&#128348; <strong>Время: %s</strong>
	&#128313; <strong>Адрес:</strong> %s
	&#128313; <strong>Герои:</strong> %s
	&#128313; <strong>Имя клиента:</strong> %s
	&#128313; <strong>Имя ребенка:</strong> %s
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
		heroes,
		fmt.Sprintf("%s %s", leadMy.Client.FirstName, leadMy.Client.LastName),
		leadMy.Child.Name,
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
		Command: "get_store_order",
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
