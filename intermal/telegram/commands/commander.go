package commands

import (
	"github.com/DaniilShd/RichShowPlatforma/intermal/config"
	"github.com/DaniilShd/RichShowPlatforma/intermal/driver"
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/repository"
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/repository/dbrepo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Commander struct {
	App *config.AppConfig
	bot *tgbotapi.BotAPI
	DB  repository.DatabaseRepoTelegram
}

func NewCommander(app *config.AppConfig, bot *tgbotapi.BotAPI, db *driver.DB) *Commander {

	return &Commander{
		App: app,
		bot: bot,
		DB:  dbrepo.NewPostgresRepoTelegram(db.SQL),
	}
}
