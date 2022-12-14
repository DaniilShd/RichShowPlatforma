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
		//Fetch

		mux.Use(Auth)
	})

	//admin pages
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(AuthAdmin)

		mux.Get("/", handlers.Repo.Dashboard)
		mux.Get("/dashboard", handlers.Repo.Dashboard)

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
		//All leads
		mux.Get("/leads-raw", handlers.Repo.AllRawLead)
		mux.Get("/leads-confirmed", handlers.Repo.AllConfirmedLead)
		mux.Get("/leads-archive", handlers.Repo.AllArchiveLead)
		mux.Get("/show-lead/{src}/{active}/{id}", handlers.Repo.ShowLead)
		//Change lead
		mux.Get("/lead-change/{src}/{active}/{id}", handlers.Repo.ChangeLead)
		mux.Post("/lead-change/{src}/{active}/{id}", handlers.Repo.ChangePostLead)
		// mux.Get("/manager", handlers.Repo.StoreItemAll)
		mux.Get("/lead-new", handlers.Repo.NewLead)
		mux.Post("/lead-new", handlers.Repo.NewPostLead)
		//Set and delete confirmed lead
		mux.Get("/lead-confirmed/{src}/{active}/{id}", handlers.Repo.SetConfirmedLead)
		mux.Get("/lead-delete-confirmed/{src}/{active}/{id}", handlers.Repo.DeleteConfirmedLead)
		//Delete lead
		mux.Get("/lead-delete/{src}/{id}", handlers.Repo.DeleteLead)

		//Leads calendar------------------------------------------------------------------------------------------
		// mux.Get("/leads-calendar", handlers.Repo.LeadsCalendar)

		//Animators
		mux.Get("/animators", handlers.Repo.AnimatorsAll)
		//New animator
		mux.Get("/animator-new", handlers.Repo.NewAnimator)
		mux.Post("/animator-new", handlers.Repo.NewPostAnimator)
		//Show and delete
		mux.Get("/animator-delete/{id}", handlers.Repo.DeleteAnimator)
		//Show anumator
		mux.Get("/animator/{id}", handlers.Repo.ShowAnimator)
		//Show change
		mux.Get("/animator-change/{id}", handlers.Repo.ChangeAnimator)
		mux.Post("/animator-change/{id}", handlers.Repo.ChangePostAnimator)

		mux.Get("/send-mail-animator", handlers.Repo.Dashboard)

		//Assistants
		mux.Get("/assistants", handlers.Repo.AssistantsAll)
		//New Assistants
		mux.Get("/assistant-new", handlers.Repo.NewAssistant)
		mux.Post("/assistant-new", handlers.Repo.NewPostAssistant)
		//Show and delete
		mux.Get("/assistant-delete/{id}", handlers.Repo.DeleteAssistant)
		//Show Assistant
		mux.Get("/assistant/{id}", handlers.Repo.ShowAssistant)
		//Show change Assistant
		mux.Get("/assistant-change/{id}", handlers.Repo.ChangeAssistant)
		mux.Post("/assistant-change/{id}", handlers.Repo.ChangePostAssistant)

		//Heroes
		mux.Get("/heroes", handlers.Repo.HeroesAll)
		//New Assistants
		mux.Get("/hero-new", handlers.Repo.NewHero)
		mux.Post("/hero-new", handlers.Repo.NewPostHero)
		//Show and delete
		mux.Get("/hero-delete/{id}", handlers.Repo.DeleteHero)
		//Show Assistant
		mux.Get("/hero/{id}", handlers.Repo.ShowHero)
		//Show change Assistant
		mux.Get("/hero-change/{id}", handlers.Repo.ChangeHero)
		mux.Post("/hero-change/{id}", handlers.Repo.ChangePostHero)

		//Fetch
		mux.Get("/fetch-leads", handlers.Repo.FetchLead)

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

		//Fetch
		mux.Get("/fetch-leads", handlers.Repo.FetchLead)
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
