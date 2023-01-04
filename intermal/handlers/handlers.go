package handlers

import (
	"net/http"

	"github.com/DaniilShd/RichShowPlatforma/intermal/config"
	"github.com/DaniilShd/RichShowPlatforma/intermal/driver"
	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
	"github.com/DaniilShd/RichShowPlatforma/intermal/render"
	"github.com/DaniilShd/RichShowPlatforma/intermal/repository"
	"github.com/DaniilShd/RichShowPlatforma/intermal/repository/dbrepo"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepository(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Dashboard(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "admin-dashboard.page.html", &models.TemplateData{})
}

func convertNumberPhone(number string) string {
	if len(number) == 10 {
		number = "+7 (" + number[0:3] + ") " + number[3:6] + "-" + number[6:8] + "-" + number[8:10]
	}
	return number
}
