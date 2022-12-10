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

		mux.Get("/", handlers.Repo.Dashboard)
		mux.Get("/dashboard", handlers.Repo.Dashboard)

		// mux.Get("/test", handlers.Repo.TestFetch)

		//Check-lists-----------------------------------------------------------------------------------------------------------------------------------
		mux.Get("/check-lists/{src}", handlers.Repo.CheckListAll)
		//Create new check-list
		mux.Get("/check-list-new/{src}", handlers.Repo.NewCheckList)
		mux.Post("/check-list-new/{src}", handlers.Repo.NewPostCheckList)
		//Changes master-class
		mux.Get("/check-lists/{src}/{id}", handlers.Repo.ShowCheckList)
		mux.Post("/check-lists/{src}/{id}", handlers.Repo.ShowPostCheckList)
		//delete any check-list
		mux.Get("/delete-check-list/{src}/{id}", handlers.Repo.DeleteCheсkList)

		//Store-----------------------------------------------------------------------------------------------------------------------------------
		mux.Get("/store", handlers.Repo.StoreItemAll)
		//Create new check-list
		mux.Get("/store-new", handlers.Repo.NewStoreItem)
		mux.Post("/store-new", handlers.Repo.NewPostStoreItem)
		//Changes master-class
		mux.Get("/store/{id}", handlers.Repo.ShowStoreItem)
		mux.Post("/store/{id}", handlers.Repo.ShowPostStoreItem)
		//delete any check-list
		mux.Get("/delete-store/{id}", handlers.Repo.DeleteStoreItem)

		//Manager
		// mux.Get("/manager", handlers.Repo.StoreItemAll)
		mux.Get("/manager-new", handlers.Repo.NewLead)

		//Leads calendar------------------------------------------------------------------------------------------
		// mux.Get("/leads-calendar", handlers.Repo.LeadsCalendar)

		//Animators
		mux.Get("/animator", handlers.Repo.Dashboard)
		mux.Get("/animator-new", handlers.Repo.Dashboard)
		mux.Get("/send-mail-animator", handlers.Repo.Dashboard)

	})

	//Stocker pages
	mux.Route("/store", func(mux chi.Router) {
		mux.Use(AuthStore)

		mux.Get("/", handlers.Repo.StoreItemAll)
		//Create new check-list
		mux.Get("/new", handlers.Repo.NewStoreItem)
		mux.Post("/new", handlers.Repo.NewPostStoreItem)
		//Changes master-class
		mux.Get("/item/{id}", handlers.Repo.ShowStoreItem)
		mux.Post("/item/{id}", handlers.Repo.ShowPostStoreItem)
		//delete any check-list
		mux.Get("/delete/{id}", handlers.Repo.DeleteStoreItem)
	})

	//Manager pages
	mux.Route("/manager", func(mux chi.Router) {
		mux.Use(AuthManager)

		// mux.Get("/", handlers.Repo.HomeManager)
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
