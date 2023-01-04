package commands

import (
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/constant"
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var startKeyBoard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButtonContact("Зарегистрироваться"),
	),
)

var AdminKeyBoard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Узнать новое"),
		tgbotapi.NewKeyboardButton("Список моих заказов"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Отправить номер телефона"),
	),
)

var ArtistKeyBoard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Мои заказы"),
		tgbotapi.NewKeyboardButton("Заказы сегодня"),
	),
)

var AssistantKeyBoard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Мои заказы"),
		tgbotapi.NewKeyboardButton("Заказы сегодня"),
	),
)

var ManagerKeyBoard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Заказы в ближайшие 2 дня"),
		tgbotapi.NewKeyboardButton("Список клиентов для обзвона"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Список артистов"),
		tgbotapi.NewKeyboardButton("Список ассистентов"),
	),
)

var StoreKeyBoard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Заказы в ближайшие 2 дня"),
		tgbotapi.NewKeyboardButton("Собранный реквизит"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Реквизит для разбора"),
		tgbotapi.NewKeyboardButton("Список всех заказов"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Список расходников"),
		tgbotapi.NewKeyboardButton("Список всех заказов"),
	),
)

func (c *Commander) Start(inputMesage *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(inputMesage.Chat.ID, "Главное меню ниже")
	switch service.ValidationChatID(inputMesage.Chat.ID) {
	case constant.ADMIN:
		msg.ReplyMarkup = AdminKeyBoard
		c.bot.Send(msg)
	case constant.MANAGER:
		msg.ReplyMarkup = ManagerKeyBoard
		c.bot.Send(msg)
	case constant.STORE:
		msg.ReplyMarkup = StoreKeyBoard
		c.bot.Send(msg)
	case constant.ARTIST:
		msg.ReplyMarkup = ArtistKeyBoard
		c.bot.Send(msg)
	case constant.ASSISTANT:
		msg.ReplyMarkup = AssistantKeyBoard
		c.bot.Send(msg)
	default:
		msg := tgbotapi.NewMessage(inputMesage.Chat.ID, "Прошу зарегистрироваться, отправив свой номер телефона.")
		msg.ReplyMarkup = startKeyBoard
		c.bot.Send(msg)
	}
}
