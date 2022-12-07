package handlers

import (
	"net/http"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
	"github.com/DaniilShd/RichShowPlatforma/intermal/render"
)

func (m *Repository) HomeStore(w http.ResponseWriter, r *http.Request) {
	pointMenu := "store"

	render.Template(w, r, "store.page.html", &models.TemplateData{
		PointMenu: pointMenu,
	})
}
