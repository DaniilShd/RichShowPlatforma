package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/constant"
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/models"
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Действие по умолчанию
func (c *Commander) DefaultBehavior(inputMesage *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(inputMesage.Chat.ID,
		"Бот умеет общаться только через меню")
	c.bot.Send(msg)
	c.Start(inputMesage)
}

// Обработчик событий для бота
func (c *Commander) HandelUpdate(update tgbotapi.Update) {

	defer func() {
		valPanic := recover()
		if valPanic != nil {
			log.Printf("panic - %v", valPanic)
		}
	}()

	if update.CallbackQuery != nil {
		var request models.RequestFromChat
		err := json.Unmarshal([]byte(update.CallbackQuery.Data), &request)
		if err != nil {
			log.Println(err)
		}

		fmt.Println(request)

		switch service.ValidationChatID(request.ChatID) {
		case constant.ADMIN:

		case constant.MANAGER:

		case constant.STORE:

		case constant.ARTIST:
			if request.Command != "" {
				switch request.Command {
				case "get_lead":
					c.GetLeadByID(request)
				default:
				}
			}
		case constant.ASSISTANT:

		default:

		}

		// msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		// 	fmt.Sprintf("Parsed %+v", request))
		// c.bot.Send(msg)
	}

	//Если присылаются данные телефона то они проверяются кому принадлежат и регистрируется чат
	if update.Message.Contact != nil {
		number := update.Message.Contact.PhoneNumber
		number = strings.ReplaceAll(number, "+", "")
		number = number[1:]
		fmt.Println(number)
		number = strings.ReplaceAll(number, ")", "")
		number = strings.ReplaceAll(number, "-", "")
		number = strings.ReplaceAll(number, "(", "")
		roleServiceConst, idAccount := service.ValidationPhoneNumber(number)
		if roleServiceConst == 0 || idAccount == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вас нет в базе данных, обратитесь к администратору Рич Шоу")
			c.bot.Send(msg)
			return
		}

		err := c.DB.SetChatIDByRoleAndID(roleServiceConst, idAccount, update.Message.Chat.ID)
		if err != nil {
			log.Fatal(err)
			return
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы успешно зарегистрировались")
		c.App.UpdateCacheAccount <- true

		c.bot.Send(msg)

		c.Start(update.Message)
		return
	}

	if update.Message != nil { // If we got a message

		switch service.ValidationChatID(update.Message.Chat.ID) {
		case constant.ADMIN:

		case constant.MANAGER:

		case constant.STORE:

		case constant.ARTIST:
			if update.Message.Command() == "" {
				switch update.Message.Text {
				case "Мои заказы":
					c.ArtistLeadsList(update.Message)
				case "Заказы сегодня":
					c.ArtistLeadsToday(update.Message)
				default:
					c.DefaultBehavior(update.Message)
				}
			} else {
				switch update.Message.Command() {
				case "start":
					c.Start(update.Message)
				default:
					c.DefaultBehavior(update.Message)
				}
			}
		case constant.ASSISTANT:

		default:
			c.Start(update.Message)
		}
	}
}
