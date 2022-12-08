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
