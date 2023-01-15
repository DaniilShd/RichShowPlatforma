package commands

import (
	"encoding/json"
	"fmt"
	"strconv"

	modelsFromApp "github.com/DaniilShd/RichShowPlatforma/intermal/models"
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *Commander) StoreOrderNew(inputMesage *tgbotapi.Message) {
	var request models.RequestFromChat

	ResponseFromApp := make(chan interface{})
	request.ResponseLeadFromApp = ResponseFromApp

	request.Command = "get_new_order"

	c.App.RequestFromTelegram <- request

	input := <-request.ResponseLeadFromApp

	storeOrders := input.([]modelsFromApp.StoreLead)

	c.bot.Send(sendMessage(inputMesage, storeOrders))

}

func sendMessage(inputMesage *tgbotapi.Message, storeOrders []modelsFromApp.StoreLead) tgbotapi.MessageConfig {
	if len(storeOrders) > 0 {
		keyBoard := tgbotapi.NewInlineKeyboardMarkup().InlineKeyboard

		count := 0
		number := 1
		var listButtonArray []tgbotapi.InlineKeyboardButton
		var res string

		for i, storeOrder := range storeOrders {
			data, _ := json.Marshal(models.RequestFromChat{
				Command: "get_order",
				ChatID:  inputMesage.Chat.ID,
				LeadID:  storeOrder.ID,
			})

			res = res + fmt.Sprintf("%d) № - %d, Дата: <strong>%s</strong>, Время: <strong>%s</strong>, Название: %s",
				number,
				storeOrder.LeadID,
				string(storeOrder.Date.Format("02-01-2006")),
				string(storeOrder.Time.Format("15:04")),
				storeOrder.Name) + "\n"

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
			if i == len(storeOrders)-1 {
				keyBoard = append(keyBoard, listButtonArray)
			}
			number++
		}
		msg := tgbotapi.NewMessage(inputMesage.Chat.ID, res)
		msg.ParseMode = "html"

		msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: keyBoard,
		}
		return msg
	} else {
		msg := tgbotapi.NewMessage(inputMesage.Chat.ID, "Заказов нет")
		return msg
	}
}

func (c *Commander) StoreOrderBuild(inputMesage *tgbotapi.Message) {
	var request models.RequestFromChat

	ResponseFromApp := make(chan interface{})
	request.ResponseLeadFromApp = ResponseFromApp

	request.Command = "get_compl_order"

	c.App.RequestFromTelegram <- request

	input := <-request.ResponseLeadFromApp

	storeOrders := input.([]modelsFromApp.StoreLead)

	c.bot.Send(sendMessage(inputMesage, storeOrders))
}

func (c *Commander) StoreOrderDestroy(inputMesage *tgbotapi.Message) {
	var request models.RequestFromChat

	ResponseFromApp := make(chan interface{})
	request.ResponseLeadFromApp = ResponseFromApp

	request.Command = "get_dest_order"

	c.App.RequestFromTelegram <- request

	input := <-request.ResponseLeadFromApp

	storeOrders := input.([]modelsFromApp.StoreLead)

	c.bot.Send(sendMessage(inputMesage, storeOrders))
}

func (c *Commander) GetStoreOrderByID(request models.RequestFromChat) {
	ResponseFromApp := make(chan interface{})
	request.ResponseLeadFromApp = ResponseFromApp

	request.Command = "get_order"

	c.App.RequestFromTelegram <- request

	input := <-request.ResponseLeadFromApp

	storeOrder := input.(modelsFromApp.StoreLead)

	var CheckListRequest models.RequestFromChat
	CheckListRequest.Command = "get_checklist"

	CheckListRequest.CheckListID = storeOrder.CheckListID

	ResponseFromAppCheck := make(chan interface{})
	CheckListRequest.ResponseLeadFromApp = ResponseFromAppCheck

	c.App.RequestFromTelegram <- CheckListRequest
	checkList := <-CheckListRequest.ResponseLeadFromApp

	checkListRes := checkList.(modelsFromApp.CheckList)

	text := "<strong>Номер заказа</strong> - " + strconv.Itoa(storeOrder.LeadID) + "\n"
	text = text + "<strong>Название программы</strong> - " + checkListRes.Name + "\n"
	text = text + "<strong>Список реквизита</strong>" + "\n"

	for _, item := range checkListRes.NameOfPoints {
		text = text + "&#9642; " + item + "\n"
	}
	text = text + "<strong>Список расходников</strong>" + "\n"

	for _, item := range checkListRes.Items {
		text = text + "&#9642; " + item.Name + " - " + strconv.Itoa(int(item.AmountItemOnce*float64(storeOrder.AmountOfChilds))) + " " + item.Dimension + "\n"
	}

	if storeOrder.Completed {

		file := tgbotapi.FilePath("./static/img/store-leads/" + storeOrder.Photo)
		res := tgbotapi.NewPhoto(request.ChatID, file)

		text = text + "<strong>Комментарий:</strong> " + storeOrder.StoreDescription

		res.ParseMode = "html"

		res.Caption = text
		c.bot.Send(res)
	} else {

		msg := tgbotapi.NewMessage(request.ChatID, text)
		msg.ParseMode = "html"

		// keyBoard := tgbotapi.NewInlineKeyboardMarkup().InlineKeyboard
		// var listButtonArray []tgbotapi.InlineKeyboardButton

		// listButton := tgbotapi.NewInlineQueryResultPhoto("12", "https://api.telegram.org/file/bot5986026405:AAEv2cbMkgQ4xNzJ60rVnrfjbB7RVuutYGE/photos/file_0.jpg")
		// listButtonArray = append(listButtonArray, listButton)

		// keyBoard = append(keyBoard, listButtonArray)

		// msg.ReplyMarkup = tgbotapi.InlineQueryResultPhoto{
		// 	Type: "jpg",
		// 	ID:   "AgACAgIAAxkBAAIFemO4duUpRy8ozdW0X-Am5eH1BKudAAJmwjEbQNDJSfbs3b6BjQGAAQADAgADeAADLQQ",

		// 	URL: "https://api.telegram.org/file/bot5986026405:AAEv2cbMkgQ4xNzJ60rVnrfjbB7RVuutYGE/photos/file_0.jpg",
		// }

		// msg.ReplyMarkup = tgbotapi.

		c.bot.Send(msg)
	}

}
