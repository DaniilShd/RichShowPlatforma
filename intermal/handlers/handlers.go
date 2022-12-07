package handlers

import (
	"github.com/DaniilShd/RichShowPlatforma/intermal/config"
	"github.com/DaniilShd/RichShowPlatforma/intermal/driver"
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
