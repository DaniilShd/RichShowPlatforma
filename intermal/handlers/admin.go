package handlers

import (
	"net/http"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
	"github.com/DaniilShd/RichShowPlatforma/intermal/render"
)

func (m *Repository) HomeManager(w http.ResponseWriter, r *http.Request) {
	pointMenu := "manager"

	render.Template(w, r, "manager.page.html", &models.TemplateData{
		PointMenu: pointMenu,
	})
}

func (m *Repository) Doashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.html", &models.TemplateData{})
}

func (m *Repository) NewMasterClass(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "admin-new-master-class.page.html", &models.TemplateData{})
}

func (m *Repository) AllMasterClass(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "admin-all-master-class.page.html", &models.TemplateData{})
}

func (m *Repository) LeadsCalendar(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "admin-leads-calendar.page.html", &models.TemplateData{})
}
