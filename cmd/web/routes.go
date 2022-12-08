package main

import (
	"net/http"

	"github.com/DaniilShd/RichShowPlatforma/intermal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func page404(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func routes() *chi.Mux {
	mux := chi.NewRouter()

	//Middleware
	mux.Use(middleware.Recoverer)
	mux.Use(WriteConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// main page
	mux.Route("/", func(mux chi.Router) {
		mux.Use(Auth)
	})

	//admin pages
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(AuthAdmin)

		mux.Get("/", handlers.Repo.HomeManager)
		mux.Get("/dashboard", handlers.Repo.Doashboard)

		//Master-class---------------------------------------------------------------------------------------------
		mux.Get("/master-class-all", handlers.Repo.AllMasterClass)
		//Changes master-class
		mux.Get("/master-class/{src}/{id}", handlers.Repo.ShowMasterClass)
		mux.Post("/master-class/{src}/{id}", handlers.Repo.ShowPostMasterClass)
		//Add new master-class
		mux.Get("/master-class-new", handlers.Repo.NewMasterClass)
		mux.Post("/master-class-new", handlers.Repo.NewPostMasterClass)
		//Delete master-class
		mux.Get("/delete-master-class/{id}", handlers.Repo.DeleteMasterClass)

		//Check-lists--------------------------------------------------------------------------------------------
		mux.Get("/check-lists/{src}/all", handlers.Repo.CheckListAll)
		//Create new check-list
		mux.Get("/check-lists-new/{src}", handlers.Repo.NewCheckList)
		mux.Post("/check-lists-new/{src}", handlers.Repo.NewPostCheckList)
		//delete any check-list
		mux.Get("/delete-check-list/{src}/{id}", handlers.Repo.DeleteCheсkList)

		//Leads calendar------------------------------------------------------------------------------------------
		mux.Get("/leads-calendar", handlers.Repo.LeadsCalendar)

		//Animators
		mux.Get("/animator-all", handlers.Repo.Doashboard)
		mux.Get("/animator-new", handlers.Repo.Doashboard)
		mux.Get("/send-mail-animator", handlers.Repo.Doashboard)

	})

	//Stocker pages
	mux.Route("/store", func(mux chi.Router) {
		mux.Use(AuthStore)

		mux.Get("/", handlers.Repo.HomeStore)
	})

	//Manager pages
	mux.Route("/manager", func(mux chi.Router) {
		mux.Use(AuthManager)

		mux.Get("/", handlers.Repo.HomeManager)
	})

	// login and logout
	mux.Get("/login", handlers.Repo.ShowLogin)
	mux.Post("/login", handlers.Repo.PostShowLogin)
	mux.Get("/logout", handlers.Repo.Logout)

	//page 404
	mux.NotFound(page404)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
