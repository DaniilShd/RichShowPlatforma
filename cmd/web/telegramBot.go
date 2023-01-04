package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/driver"
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/commands"
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartTelegramBot() {
	bot, err := tgbotapi.NewBotAPI("5986026405:AAEv2cbMkgQ4xNzJ60rVnrfjbB7RVuutYGE")
	if err != nil {
		fmt.Println(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.UpdateConfig{
		Timeout: 600,
	}

	log.Println("Connecting to database...")
	dsn := fmt.Sprintf("host=localhost port=5432 dbname=richshow user=postgres password=root")
	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.SQL.Close()
	log.Println("Connected to database!")

	// file := tgbotapi.FilePath("./static/img/rich_show.png")
	// f := tgbotapi.NewPhoto(379017783, file)

	// fmt.Println(f)
	// f.Caption = "Hello"
	// bot.Send(f)
	// bot.Send(msg)

	updates := bot.GetUpdatesChan(u)

	// accounts := accounts.NewAccounts(db)

	newCommander := commands.NewCommander(&app, bot, db)

	//Создаем мапу аккаунтов

	//Горутина для обновления аккаунтов, если зарегистируется кто то новый, через менеджера или админа
	ctx := context.Background()
	go updateAccounts(ctx, db)

	//Передаем в канал любые сообщения, в любой момент времени, они сразу отправляются адресату
	go listenChannelMail(bot)

	//Напоминалка о праздниках
	go orderRemind(bot, db)

	//читаем входящие свообщения для бота
	for update := range updates {
		newCommander.HandelUpdate(update)
	}

}

func listenChannelMail(bot *tgbotapi.BotAPI) {
	for {
		data := <-app.MailChan
		if data.ChatID != 0 {
			message := tgbotapi.NewMessage(data.ChatID, data.Text)
			bot.Send(message)
		}
	}
}

func updateAccounts(ctx context.Context, db *driver.DB) {
	store.NewAccount(db)
	for {
		timer := time.NewTicker(time.Second * 2)
		defer timer.Stop()

		select {
		case <-timer.C:
			store.NewAccount(db)
		case <-app.UpdateCacheAccount:
			store.NewAccount(db)

		case <-ctx.Done():
			return
		}
	}
}

func orderRemind(bot *tgbotapi.BotAPI, db *driver.DB) {
	//Дописать код по оповещению сотрудников о праздниках за сутки и за 2 часа, менеджер подтверждает заказы за 2 дня
}
