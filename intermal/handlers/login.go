package handlers

import (
	"net/http"

	"github.com/DaniilShd/RichShowPlatforma/intermal/forms"
	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
	"github.com/DaniilShd/RichShowPlatforma/intermal/render"
	modelsTelegram "github.com/DaniilShd/RichShowPlatforma/intermal/telegram/models"
)

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	m.App.MailChan <- modelsTelegram.MailData{
		ChatID: 379017783,
		Text:   "Вход в систему",
	}

	var login string
	var password string
	login = r.Form.Get("login")
	password = r.Form.Get("password")

	stringMap := make(map[string]string)
	stringMap["login"] = login
	stringMap["password"] = password

	form := forms.New(r.PostForm)
	form.Required("login", "password")

	if !form.Valid() {
		m.App.Session.Put(r.Context(), "error", "Некорректно заполнены поля")
		render.Template(w, r, "login.page.html", &models.TemplateData{
			Form:      form,
			StringMap: stringMap,
		})
		return
	}

	id, access_level, _, err := m.DB.Authenticate(login, password)
	if err != nil {
		// log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Неверный логин или пароль")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "access_level", access_level)
	m.App.Session.Put(r.Context(), "flash", "Успешный вход")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
