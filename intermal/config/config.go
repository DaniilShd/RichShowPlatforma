package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// Appconfid holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Session       *scs.SessionManager
	InProduction  bool
}
