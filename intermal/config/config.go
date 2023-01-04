package config

import (
	"html/template"
	"log"

	modelsTelegram "github.com/DaniilShd/RichShowPlatforma/intermal/telegram/models"
	"github.com/alexedwards/scs/v2"
)

// Appconfid holds the application config
type AppConfig struct {
	UseCache            bool
	TemplateCache       map[string]*template.Template
	InfoLog             *log.Logger
	ErrorLog            *log.Logger
	Session             *scs.SessionManager
	InProduction        bool
	MailChan            chan modelsTelegram.MailData
	UpdateCacheAccount  chan bool
	RequestFromTelegram chan modelsTelegram.RequestFromChat
}
