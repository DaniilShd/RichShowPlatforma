package models

import "github.com/DaniilShd/RichShowPlatforma/intermal/forms"

type TemplateData struct {
	StringMap      map[string]string
	IntMap         map[string]int
	FloatMap       map[string]float64
	Data           map[string]interface{}
	CSRFToken      string
	Flash          string
	Warning        string
	Error          string
	Form           *forms.Form
	IsAuthenticate int
	AccessLevel    int
	PointMenu      string
}
